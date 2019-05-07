# MTS Communicator mock
This service is mock/stub/fixture for [MTS Communicator](https://api.mcommunicator.ru/m2m/m2m_api.asmx?wsdl). It can be usefull for testing interaction with the MTS communicator.

The service allows to receive and send test SMS, and provide simple UI for that.

For the moment the following methods are implemented (not fully):
+ `GetMessages`
+ `SendMessage`
+ `GetMessagesStatus`

There are some UI pages for testing:
+ List outgoing messages (from service) is located on `/ui/list`. The current implementation shows only last 10 messages. In future it will be available to add filters and use pagination.
+ UI for create incoming message (as user) is located on `ui/send`.

The service listening `9000` port on all interfaces and provide SOAP intaraction by `/test.svc` url.
