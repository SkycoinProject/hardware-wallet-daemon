# Hardware Wallet Daemon API
API default service port is `9510`.

The API currently supports skywallet and its emulator.

The skywallet endpoints start with `/api` and emulator endpoints with `/api/emulator`.

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

- [Usage](#usage)
    - [Main Endpoints](#main-endpoints)
        - [Generate Addresses](#generate-addresses)
        - [Apply Settings](#apply-settings)
        - [Backup Seed](#backup-seed)
        - [Cancel](#cancel)
        - [Check Message Signature](#check-message-signature)
        - [Get Features](#get-features)
        - [Firmware Update](#firmware-update)
        - [Recover Wallet](#recover-old-wallet)
        - [Generate Mnemonic](#generate-mnemonic)
        - [Set Mnemonic](#set-mnemonic)
        - [Set Pin Code](#set-pin-code)
        - [Sign Message](#sign-message)
        - [Transaction Sign](#transaction-sign)
        - [Wipe](#wipe)
        
    

<!-- /MarkdownTOC -->
## Main Endpoints

### Generate Addresses
Generate addresses for the hardware wallet seed.

```
URI: /api/v1/generateAddresses
Method: POST
Content-Type: application/json
Args: {"address_n": "<address_n>", "start_index": "<start_index>", "confirm_address": "<confirm_address>"}
```

Example:
```sh
$ curl http://127.0.0.1:6430/api/v1/generateAddresses \
  -H 'Content-Type: application/json' \
  -d '{"address_n": 2, "start_index": 0}'
```

Response:
```json
{
    "data": {
        "addresses": [
            "XA3kV3QYWF9QnktordpYggujDRcGXg4tpQ",
            "2ew81SqY1C5efFA2ULQ4NvJNhxWnM1jPGxu"
        ]
    }
}
```

### Apply Settings
Apply hardware wallet settings.

```
URI: /api/v1/applySettings
Method: POST
Args:
    label: label for hardware wallet
    use-passphrase: (boolean) ask for passphrase before starting operation
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/applySettings \
   -H 'Content-Type: application/x-www-form-urlencoded' \
   -d 'label=skywallet'
```

Response:
```json
{
    "data": "Settings applied"
}
```

### Backup Seed
Start seed backup procedure

```
URI: /api/v1/backup
Method: POST
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/backup
```

Response Flow:
- Button confirmation requests are shown for each word of the seed on the hardware wallet screen.
- The whole procedure is repeated again.
- Success response

```json
{
    "data": "Device backed up!"
}
```

- If device has no seed
```json
{
    "error": {
        "message": "Device not initialized",
        "code": 409
    }
}
```

### Cancel
Cancels the current operation.

> This function can be called safely even if no operation is active. The response will be the same.

```
URI: /api/v1/cancel
Method: PUT
```

Example:
```sh
$ curl -X PUT http://127.0.0.1:6430/api/v1/cancel
```

Response:
```json
{
    "data": "Action cancelled by user"
}
```

### Check Message Signature
Check a message signature matches the given address.

```
URI: /api/v1/checkMessageSignature
Method: POST
Content-Type: application/json
Args: {"message": "<message>", "signature": "<signature>", "address": "<address>"}
```

Example:
```sh
curl -X POST http://127.0.0.1:6430/api/v1/checkMessageSignature \
-H 'Content-Type: application/json' \
-d '{"message": "Hello World!", "signature": "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn", "address": "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"}'
```

### Get Features
Returns device information.

```
URI: /api/v1/features
Method: GET
```

Example:
```sh
$ curl http://127.0.0.1:6430/api/v1/features
```

Response:
```json
{
    "data": {
        "features": {
            "vendor": "Skycoin Foundation",
            "major_version": 1,
            "minor_version": 6,
            "patch_version": 1,
            "device_id": "92FD1C8A06E401ABC524DF1C",
            "pin_protection": false,
            "passphrase_protection": false,
            "label": "skywallet",
            "initialized": true,
            "bootloader_hash": "/V5yWwVclKtbpzZW/M1CeuNxQZQZTv5ia6OIOiVu5C8=",
            "pin_cached": false,
            "passphrase_cached": false,
            "needs_backup": false,
            "model": "1"
        }
    }
}
```


### Firmware Update
Update device firmware

```
URI: /api/v1/firmwareUpdate
Method: PUT
Args:
    file: firmware file
```


### Recover Wallet
Recover existing wallet using seed.

The device needs to be wiped if already initialized.

```
URI: /api/v1/recovery
Method: POST
Args:
    word-count: mnemonic seed length
    use-passphrase: (boolean) ask for passphrase before starting operation
    dry-run: (bool) perform dry-run recovery workflow (for safe mnemonic validation).
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/recovery \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'word-count=$word-count' \
  -d 'use-passphrase=$use-passphrase' \
  -d 'dry-run=$dry-run'
```

Response Flow:
- The user is asked for confirmation on hardware wallet screen to start the recovery process.
- The daemon then returns intermediate word request type in response.
- The frontend needs to send intermediate word request till the hardware wallet keeps showing instructions on screen.
- Sucess response:

```json
{
    "data": "Device recovered"
}
```

### Generate Mnemonic
Generate mnemonic can be used to initialize the device with a random seed.

```
URI: /api/v1/generateMnemonic
Method: POST
Args:
    word-count: mnemonic seed length
    use-passphrase: (bool) ask for passphrase before starting operation
```

Example:
```sh
$ curl http://127.0.0.1:6430/api/v1/generateMnemonic \
  -H 'Content-Type: x-www-form-urlencoded' 
  -d 'word-count=12' \
  -d 'use-passphrase=$use-passphrase'
```

Response:
```json
{
    "data": "Mnemonic successfully configured"
}
```

### Set Mnemonic
Set mnemonic can be used to initialize the device with your own seed.

> The seed needs to be a valid bip39 seed.

```
URI: /api/v1/setMnemonic
Method: POST
Args:
    mnemonic: bip39 mnemonic seed [required]
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/setMnemonic\
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'mnemonic=$mnemonic'
```

Response:
- Valid mnemonic

```json
{
    "data": "cloud flower upset remain green metal below cup stem infant art thank"
}
```

- Invalid mnemonic
```json
{
    "error": {
        "message": "Mnemonic with wrong checksum provided",
        "code": 409
    }
}
```

### Set Pin Code
Configure a pin code on the device.

```
URI: /api/v1/setPinCode
Method: POST
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/setPinCode
```

Response Flow:
- User is shown a button confirmation request on hardware wallet to confirm start of pin code process.
- The daemon returns intermediate pinmatrix request two times to the frontend.
- Response on success

```json
{
    "data": "PIN changed"
}
```
- Response on failure
```json
{
    "error": {
        "message": "PIN mismatch",
        "code": 409
    }
}
```

### Sign Message
Sign a message using the secret key at given index.

```
URI: /api/v1/signMessage
Method: POST
Args: {
    "address_n": <address_n>, 
    "message": "<message>"
}
```

Example:
```sh
$ curl -X POST http://127.0.0.1:6430/api/v1/signMessage \
  -H 'Content-Type: application/json' \
  -d '{"address_n": 0, "message": "hello world"}'
```

Response:
```json
{
    "data": {
        "signature": "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQWBXZ4WospXL7bJNu7aSVn5eCPrATSkGjtQfzGYMNFQDYt"
    }
}
```

### Transaction Sign
Sign a transaction with the hardware wallet.

```
URI: /api/v1/transactionSign
Method: POST
Args: {
    "inputs": "<inputs>", 
    "input_indexes": "<input_indexes>", 
    "output_addresses": "<output_addresses>", 
    "coins": "<coins>", 
    "hours": "<hours>", 
    "address_indexes": "<address_indexes>"
    }
```

Example:
```sh
$ curl http://127.0.0.1:6430/api/v1/transactionSign \
  -H 'Content-Type: application/json' \
  -d '{"inputs": ["e3411a073376d2abf2e3231023fca48f2396c575871764276d6206350207cde4"],
    "input_indexes": [0],
    "output_addresses": ["2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"],
    "coins": ["0.5"],
    "hours": ["1"],
    "address_indexes": [1]
   	}
```

### Wipe
Wipe deletes all data from the hardware wallet.

```
URI: /api/v1/wipe
Method: DELETE
```

Example:
```sh
$ curl -X DELETE http://127.0.0.1:6430/api/v1/wipe
```

Response:
```json
{
    "data": "Device wiped"
}
```
