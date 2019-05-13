# Hardware Wallet Daemon API
API default service port is `9510`.

The API currently supports skywallet and its emulator.

The skywallet endpoints start with `/api/v1` and emulator endpoints with `/api/v1/emulator`.

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
        - [Configure Pin Code](#configure-pin-code)
        - [Sign Message](#sign-message)
        - [Transaction Sign](#transaction-sign)
        - [Wipe](#wipe)
        - [Available](#available)
    

<!-- /MarkdownTOC -->
## Main Endpoints

### Generate Addresses
Generate addresses for the hardware wallet seed.

```
URI: /api/v1/generate_addresses
Method: POST
Content-Type: application/json
Args: {"address_n": "<address_n>", "start_index": "<start_index>", "confirm_address": "<confirm_address>"}
```

**Parameters**
- `address_n`: Number of addresses to generate. Assume 1 if not set.
- `start_index`: Index where deterministic key generation will start from. Assume 0 if not set.
- `confirm_address`: If requesting one address it will be sent only if user confirms operation by pressing device's button.

Example:
```sh
$ curl http://127.0.0.1:9510/api/v1/generate_addresses \
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
URI: /api/v1/apply_settings
Method: POST
Content-Type: application/json
Args: {"label": "<label for hardware wallet>", "use_passphrase": "<ask for passphrase before starting operation>"}
```

Example:
```sh
$ curl -X POST http://127.0.0.1:9510/api/v1/apply_settings \
   -H 'Content-Type: application/json' \
   -d '{"label": "skywallet", "use_passphrase": false}'
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
$ curl -X POST http://127.0.0.1:9510/api/v1/backup
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
```bash
$ curl -X PUT http://127.0.0.1:9510/api/v1/cancel
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
URI: /api/v1/check_message_signature
Method: POST
Content-Type: application/json
Args: {"message": "<message>", "signature": "<signature>", "address": "<address>"}
```

**Parameters**
- `message`: The message that the signature claims to be signing.
- `signature`: Signature of the message.
- `address`: Address that issued the signature.

Example:
```bash
curl -X POST http://127.0.0.1:9510/api/v1/check_message_signature \
-H 'Content-Type: application/json' \
-d '{"message": "Hello World", "signature": "6ebd63dd5e57cad07b6d229e96b5d2ac7d1bec1466d2a95bd200c21be6a0bf194b5ad5123f6e37c6393ee3635b38b938fcd91bbf1327fc957849a9e5736f6e4300", "address": "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"}'
```

### Get Features
Returns device information.

```
URI: /api/v1/features
Method: GET
```

Example:
```bash
$ curl http://127.0.0.1:9510/api/v1/features
```

Response:
```json
{
    "data": {
        "vendor": "Skycoin Foundation",
        "device_id": "617F30E0E10C9C9E93B5EE37",
        "pin_protection": false,
        "passphrase_protection": false,
        "label": "617F30E0E10C9C9E93B5EE37",
        "initialized": false,
        "bootloader_hash": "k3jy1kG+JbHZxi2cTSScekQK007/YkbZzyXWI6dVJns=",
        "pin_cached": false,
        "passphrase_cached": false,
        "needs_backup": false,
        "model": "1",
        "fw_major": 1,
        "fw_minor": 7,
        "fw_patch": 0,
        "firmware_features": 4
    }
}
```


### Firmware Update
Update device firmware

```
URI: /api/v1/firmware_update
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
Content-Type: application/json
Args: {
        "word_count": "<mnemonic seed length>", 
        "use_passphrase": "<ask for passphrase before starting operation>", 
        "dry_run": "<perform dry-run recovery workflow (for safe mnemonic validation)>"
      }
```

Example:
```bash
$ curl -X POST http://127.0.0.1:9510/api/v1/recovery \
  -H 'Content-Type: application/json' \
  -d '{"word_count": 12, "use_passphrase": false, "dry_run": true}'
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
URI: /api/v1/generate_mnemonic
Method: POST
Content-Type: application/json
Args: {"word_count": "<mnemonic seed length>", "use_passphrase": "<ask for passphrase before starting operation>"}
```

Example:
```bash
$ curl http://127.0.0.1:9510/api/v1/generate_mnemonic \
  -H 'Content-Type: application/json' \
  -d '{"word_count": 12, "use_passphrase": false}'
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
URI: /api/v1/set_mnemonic
Method: POST
Content-Type: application/json
Args: {"mnemonic": "<bip39 mnemonic seed>"}
```

Example:
```bash
$ curl -X POST http://127.0.0.1:9510/api/v1/set_mnemonic \
  -H 'Content-Type: application/json' \
  -d '{"mnemonic": "cloud flower upset remain green metal below cup stem infant art thank"}'
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

### Configure Pin Code
Configure a pin code on the device.


```
URI: /api/v1/configure_pin_code
Method: POST
Args:
- remove_pin: (optional) Used to remove current pin
```

Example:
```bash
$ curl -X POST http://127.0.0.1:9510/api/v1/configure_pin_code
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

Example(Remove Pin):
```bash
$  curl -X POST http://127.0.0.1:9510/api/v1/configure_pin_code \
   -H 'Content-Type: application/x-www-form-urlencoded' \
   -d 'remove_pin=true'
```

Response Flow:
- User is shown a button confirmation request on hardware wallet to confirm start of pin remove process.
- The daemon returns intermediate pinmatrix request once to the frontend.
- Response on success

```json
{
    "data": "PIN removed"
}
```
- Response on failure
```json
{
    "error": {
        "message": "PIN invalid",
        "code": 409
    }
}
```

### Sign Message
Sign a message using the secret key at given index.

```
URI: /api/v1/sign_message
Method: POST
Args: {
    "address_n": <address_n>, 
    "message": "<message>"
}
```

**Parameters**
- `address_n`: Index of the address that will issue the signature.
- `message`: The message that the signature claims to be signing.

Example:
```bash
$ curl -X POST http://127.0.0.1:9510/api/v1/sign_message \
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
URI: /api/v1/transaction_sign
Method: POST
Args: {
    "transaction_inputs": [{"index": <index<, "hash":"<hash>"}],
    "transaction_outputs": [{"address_index": <address_index>,"address":"<address>","coins":"<coins>","hours":"<hours>"}],
   } 
```

**Parameters**
- transaction_inputs: List of objects with the following fields:
  * `index`: Index of the address, in the hardware wallet, to which the input belongs.
  * `hash`: Input hash.
- transaction_outputs: List of objects with the following fields:
  * `address_index`: If the output is used for returning coins/hours to one of the addresses of the hardware
  * `address`: Skycoin address in `Base58` format.
  wallet, this parameter must contain the index of the address in the hardware wallet, so that the user is
  not asked for confirmation for this specific output. If this is not the case, this parameter is not necessary.
  * `coins`: Output coins.
  * `hours`: Output hours.

Example:
```bash
$ curl http://127.0.0.1:9510/api/v1/transaction_sign \
  -H 'Content-Type: application/json' \
  -d '{"transaction_inputs":[{"index":0,"hash":"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},{"index":1,"hash":"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"}],"transaction_outputs":[{"address_index":null,"address":"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG","coins":"2","hours":"2"},{"address_index":null,"address":"2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8","coins":"3","hours":"3"}]}'
```

### Wipe
Wipe deletes all data from the hardware wallet.

```
URI: /api/v1/wipe
Method: DELETE
```

Example:
```bash
$ curl -X DELETE http://127.0.0.1:9510/api/v1/wipe
```

Response:
```json
{
    "data": "Device wiped"
}
```

### Available
Available tells whether a device is currently connected to the machine or not.

```
URI: /api/v1/available
Method: GET
```

Example:
```bash
$ curl -X GET http://127.0.0.1:9510/api/v1/available
```

Response:
```json
{
    "data": true
}
```

