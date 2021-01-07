# Copyright 2019 The OpenSDS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import functools
import time
import json
import threading

import log
from utils import base
from utils import config as cfg
from kafka import KafkaConsumer
import pandas as pd
import numpy as np
import tensorflow as tf

LOG = log.getLogger(__name__)
CONF = cfg.CONF


data_parser_opts = [
    cfg.StrOpt('kafka_bootstrap_servers',
               default='localhost:9092',
               help='kafka bootstrap server'),
    cfg.StrOpt('kafka_topic',
               default='delfin-kafka',
               help='kafka topic'),
    cfg.IntOpt('kafka_retry_num',
               default=3,
               help='kafka retry num')
]

CONF.register_opts(data_parser_opts, "data_parser")

DATA_ENTRIES_COUNT = 500

class DataReceiver(base.Base):
    def __init__(self, name):
        super(DataReceiver, self).__init__()
        self._name = name

    def run(self):
        raise NotImplemented


class DataDictionary:
    def __init__(self):
        self.dict = {}

    def get(self, key):
        return self.dict.get(key)

    def get(self):
        return self.dict

    def update(self, key, value):
        if key in self.dict.keys():
            # sample only DATA_ENTRIES_COUNT
            if len(self.dict[key]) >= DATA_ENTRIES_COUNT:
                self.dict[key].pop(0)
        else:
            self.dict[key] = []

        LOG.debug("Updating metric for %s with %s", key, value)
        self.dict[key].append(value)

    def has_key(self, key):
        return key in self.dict.keys()

    def print(self):
        LOG.info(str(self.dict))

    def len(self, key):
        return len(self.dict[key])


class ModelDictionary:
    def __init__(self):
        self.dict = {}

    def get(self, key):
        return self.dict.get(key)

    def update(self, key, value):
        self.dict[key] = value

    def has_key(self, key):
        return key in self.dict.keys()

    def print(self):
        LOG.info(str(self.dict))


class TrainingDictionary:
    def __init__(self):
        self.dict = {}

    def get(self, key):
        return self.dict.get(key)

    def update(self, key, value):
        self.dict[key] = value

    def add_entry(self, dict_key, key, value):
        dict_val = self.dict[dict_key]
        dict_val[key] = value
        self.dict[dict_key] = dict_val

    def has_key(self, key):
        return key in self.dict.keys()

    def print(self):
        LOG.info(str(self.dict))


data_dictionary = DataDictionary()
model_dictionary = ModelDictionary()
training_dictionary = TrainingDictionary()

TIME_STEPS = 288


# Generated training sequences for use in the model.
def create_sequences(values, time_steps=TIME_STEPS):
    output = []
    for i in range(len(values) - time_steps):
        output.append(values[i : (i + time_steps)])

    return np.stack(output)


# Visualize the data
# Timeseries data without anomalies
def create_training_value(key, values):
    train_data = {}
    for value in values:
        val_keys = list(value.keys())[0]
        train_data[val_keys] = [value[val_keys]]

    train_df = pd.DataFrame.from_dict(train_data, orient='index', columns=['read_bandwidth'])

    if training_dictionary.has_key(key):
        # if training is already available, case of checking anomaly
        train_dict = training_dictionary.get(key)
        df_training_value = (train_df - train_dict.get('Mean')) / train_dict.get('STD')
        x_train = create_sequences(df_training_value.values)
        LOG.info("Training input shape for anomaly detection: %s", x_train.shape)
        return x_train
    else:
        # case for first time training
        training_mean = train_df.mean()
        training_std = train_df.std()
        training_dictionary.update(key, {'Mean': training_mean, 'STD': training_std})
        df_training_value = (train_df - training_mean) / training_std
        x_train = create_sequences(df_training_value.values)
        LOG.info("Training input shape for training: %s", x_train.shape)
        return x_train


def create_model(x_train):
    """
    ## Build a model
    We will build a convolutional reconstruction autoencoder model. The model will
    take input of shape `(batch_size, sequence_length, num_features)` and return
    output of the same shape. In this case, `sequence_length` is 288 and
    `num_features` is 1.
    """
    model = tf.keras.Sequential(
        [
            tf.keras.layers.Input(shape=(x_train.shape[1], x_train.shape[2])),
            tf.keras.layers.Conv1D(
                filters=32, kernel_size=7, padding="same", strides=2, activation="relu"
            ),
            tf.keras.layers.Dropout(rate=0.2),
            tf.keras.layers.Conv1D(
                filters=16, kernel_size=7, padding="same", strides=2, activation="relu"
            ),
            tf.keras.layers.Conv1DTranspose(
                filters=16, kernel_size=7, padding="same", strides=2, activation="relu"
            ),
            tf.keras.layers.Dropout(rate=0.2),
            tf.keras.layers.Conv1DTranspose(
                filters=32, kernel_size=7, padding="same", strides=2, activation="relu"
            ),
            tf.keras.layers.Conv1DTranspose(filters=1, kernel_size=7, padding="same"),
        ]
    )

    model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=0.001), loss="mse")
    # model.summary()

    return model


def training_model(model, x_train):
    """
    ## Train the model
    Please note that we are using `x_train` as both the input and the target
    since this is a reconstruction model.
    """
    history = model.fit(
        x_train,
        x_train,
        epochs=50,
        batch_size=128,
        validation_split=0.1,
        callbacks=[
            tf.keras.callbacks.EarlyStopping(monitor="val_loss", patience=5, mode="min")
        ],
    )

    LOG.debug("Loss in training %s ", history.history["loss"])
    LOG.debug("Validation loss in training %s ", history.history["val_loss"])
    return model


SLEEP_SEC = 15


def process_data_dictionary():
    while True:
        LOG.info("--------------> Sleeping for %s sec <--------- ", SLEEP_SEC)
        time.sleep(SLEEP_SEC)

        data_dict = data_dictionary.get()
        for key in data_dict:
            # Storage ID is key

            LOG.info("Processing storage[%s] metrics", key)

            # Skip processing if entries less than DATA_ENTRIES_COUNT
            if data_dictionary.len(key) < DATA_ENTRIES_COUNT:
                LOG.info("Skipping processing as entries are less for storage[%s] len(%s)", key, data_dictionary.len(key))
                continue

            # Training values
            x_train = create_training_value(key, data_dict[key])

            if training_dictionary.has_key(key) and 'Threshold' in training_dictionary.get(key).keys():
                # get the prediction
                x_test_pred = model_dictionary.get(key).predict(x_train)
                test_mae_loss = np.mean(np.abs(x_test_pred - x_train), axis=1)
                test_mae_loss = test_mae_loss.reshape((-1))

                # print anomalies
                threshold = training_dictionary.get(key)['Threshold']
                anomalies = test_mae_loss > threshold

                # Check for anomaly
                anomaly = False
                for x in anomalies:
                    if x:
                        LOG.warning("Anomaly detected")
                        anomaly = True
                LOG.info("Anomaly not detected")
            else:
                model = create_model(x_train)
                train_model = training_model(model, x_train)

                # Update the trained model
                model_dictionary.update(key, train_model)

                # Get train MAE loss.
                x_train_pred = train_model.predict(x_train)
                train_mae_loss = np.mean(np.abs(x_train_pred - x_train), axis=1)

                # Get reconstruction loss threshold.
                threshold = np.max(train_mae_loss)
                training_dictionary.add_entry(key, 'Threshold', threshold)
                LOG.info("Training loss threshold : %s", threshold)



class KafkaDataReceiver(DataReceiver):
    def __init__(self):
        super(KafkaDataReceiver, self).__init__(name="kafka")

    def consume(self):
        consumer = KafkaConsumer(CONF.data_parser.kafka_topic,
                                 bootstrap_servers=CONF.data_parser.kafka_bootstrap_servers,
                                 auto_offset_reset='earliest')

        for msg in consumer:
            perf = json.loads(msg.value)
            # Extract storage_id
            storage_id = [elem for elem in perf[0][1].values()][0]
            LOG.debug("Adding metric for storage_id : %s", storage_id)
            for data in perf:
                if data[0] == 'read_bandwidth':
                    LOG.debug("Data to be updated : %s", [elem for elem in data[2].values()][0])
                    if [elem for elem in data[2].values()][0] > 100:
                        LOG.warning("Should detect anomaly for [%s]", storage_id)
                    data_dictionary.update(storage_id, data[2])
                    break

    def run(self):
        retry = CONF.data_parser.kafka_retry_num
        for index in range(1, retry+1):
            try:
                self.consume()
            except KeyboardInterrupt:
                LOG.info("Bye!")
                break
            except Exception as e:
                if index > retry:
                    LOG.error('%s\nall retry failed, exit.', e)
                    raise
                else:
                    LOG.error("%s ,retry %d time(s)", e, index)
            else:
                break


class Manager(base.Base):
    def __init__(self, receiver_name):
        super(Manager, self).__init__()
        self._receiver = KafkaDataReceiver()

    def run(self):
        try:
            thread = threading.Thread(target=process_data_dictionary)
            thread.start()
            self._receiver.run()
            thread.join()
        except Exception as e:
            LOG.error("%s ", e)

