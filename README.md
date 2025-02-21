# SMS

Send SMS messages using RDCom API.

API info at [API site](https://docs.rdcom.com/docs/platform/intro).



The application automatically loads a `.env` file if available; the format is that of the checked-in `.env`. In order for the local file not to be overwritten by the upstream version do the following:

```bash
$> git update-index --assume-unchanged .env
```