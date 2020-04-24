# golang-iot-home
Webserver written in Golang for iot-dashboard.

As a test, rewrote a web server I have written in .NET Core in Golang in order to learn Go language and it's caveats.

To show data in the application you can configure
- Netatmo weather station
- Wunderlist lists, up to 5 different lists
- Philips Hue bridge


Run 'go run main.go' to start the app. Navigate to localhost:3001 to see the app.
Port is hardcoded since this is not intended to be any production app, just a learning process.

Use the following environment variables, located in root in '.env' to make it work.
Need to press Philips Hue button first time you start the server in order to link the app to the Hue Bridge.

- NETATMO_CLIENTID: 
- NETATMO_CLIENTSECRET: 
- NETATMO_DEVICEID: 
- NETATMO_USERNAME: 
- NETATMO_PASSWORD: 
- NETATMO_OUTDOORMODULEID: 
- WUNDERLIST_ACCESSTOKEN: 
- WUNDERLIST_CLIENTID: 
- WUNDERLIST_LISTFIRST: 
- WUNDERLIST_LISTSECOND: 
- WUNDERLIST_LISTTHIRD: 
- PHILIPS_HUE_BRIDGEIP: 
