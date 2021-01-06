The example is POC of Anomaly detection of delfin metric data
Kafka is used data stream for receiving metric data which is plugged in into Anomaly detection module.

### Architecture

````
 ---------            ---------            ------------------
| Delfin  | ------>  |  Kafka  |  ----->  | Anomaly Detection |
 ---------            ---------            -------------------
````

- Delfin publish metric data on kafka (Default topic: delfin-kafka)
- Anomaly detection module subscribe for the topic and store the metric data in dictionary
- Anomaly detections runs on stored data in dictionary

### Constraints

- Currently only "read_bandwidth" considered
- (TODO) Online model training is not implemented 
- Only 1000 samples are considered for model training

### Steps:
- Setup Delfin and Kafka with help of [documentation](https://docs.sodafoundation.io/guides/developer-guides/delfin/)
- install dependencies mentioned in requirements.txt
- export PYTHONPATH=`pwd`
- python3 cmds/parser.py --config-file etc/anomaly_detection.conf


## Documentation

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Quick Start - To Use/Experience

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Quick Start - To Develop

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Latest Releases

[https://github.com/sodafoundation/anomaly-detection/releases](https://github.com/sodafoundation/anomaly-detection/releases)

## Support and Issues

[https://github.com/sodafoundation/anomaly-detection/issues](https://github.com/sodafoundation/anomaly-detection/issues)

## Project Community

[https://sodafoundation.io/slack/](https://sodafoundation.io/slack/)

## How to contribute to this project?

Join [https://sodafoundation.io/slack/](https://sodafoundation.io/slack/) and share your interest in the ‘general’ channel

Checkout [https://github.com/sodafoundation/anomaly-detection/issues](https://github.com/sodafoundation/anomaly-detection/issues) labelled with ‘good first issue’ or ‘help needed’ or ‘help wanted’ or ‘StartMyContribution’ or ‘SMC’

## Project Roadmap

Provide advanced anomaly detection and prediction based on the resource or performance data.

[https://docs.sodafoundation.io](https://docs.sodafoundation.io/)

## Join SODA Foundation

Website : [https://sodafoundation.io](https://sodafoundation.io/)

Slack  : [https://sodafoundation.io/slack/](https://sodafoundation.io/slack/)

Twitter  : [@sodafoundation](https://twitter.com/sodafoundation)

Mailinglist  : [https://lists.sodafoundation.io](https://lists.sodafoundation.io/)