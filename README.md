# OAuth2

BITGIN's OAuth implementation supports the standard [authorization code grant type](https://www.rfc-editor.org/rfc/rfc6749#section-4.1)

- [What is OAuth 2.0 ?](https://itnext.io/an-oauth-2-0-introduction-for-beginners-6e386b19f7a9)
- [The library we use](https://github.com/go-oauth2/oauth2)


> We also supply simple oauth2 client example to help you understand the authorization flow specifically and test easily. 

# Table of contents
- [Quick Start](#quick-start)
- [Authentication](#authentication)
  - [How to get access token ?](#how-to-get-access-token)
  - [How to refresh access token ?](#how-to-refresh-access-token)

- [Exchange API](#exchange-api)
  - Account
    - [Get Account](#get-account)
    - [Get Account Bank](#get-account-bank)
  - Wallet
    - [Get Balance](#get-balance)
    - [Get Deposit History](#get-deposit-history)
    - [Get Deposit Address](#get-deposit-address)
    - [Get Deposit Bank](#get-deposit-bank)
    - [Get Withdrawal History](#get-withdrawal-history)
    - [Request Withdrawal](#request-withdrawal)
    - [Confirm Withdrawal](#confirm-withdrawal)
  - Trade
    - [Get Trade History](#get-trade-history)
    - [Get Quote](#get-quote)
    - [Accept Quote](#accept-quote)
  - Appendix
    - [Withdrawal Type Definition](#withdrawal-type-definition)
    - [Withdrawal Status Definition](#withdrawal-status-definition)
    - [Deposit Status Definition](#deposit-status-definition)
    - [Order Side Definition](#order-side-definition)

## Quick Start

- Installation
    ```
    $ git clone github.com/bitgin/bitgin-oauth2-client-example
    ```

- Run executable file (arm64) or you can build executables for different architectures [here](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)

    ```
    $ oa2cli -h

    Usage: oa2cli [options] 
    Currently, the following flags can be used
    -e string
            BITGIN environment mode e.g. stage, prod (default "stage")
    -i string
            client id
    -p string
            client serve port (default "9094")
    -s string
            client secret
    -u string
            user id
    ```
- example 
    ```
    $ oa2cli -i [clientID]  -s [clientSecret] -u [targetUserID]
    ```

## Authentication

### How to get access token ?

First of all, the entry point of OAuth2 Authorization process, images you have a button on your application UI, when user trigger the button, it will be redirect to  BITGIN Domain `GET /oauth/authorize`  with the following parameters

**Header**
| Key | Value |
| --- | --- |
| Content-Type | application/json |

**Parameters**
        
| Field | Description |
| --- | --- |
| client_id | represents your client id. |
| code_challenge | the value of code_verifier after sha256 hash. |
| user_id | represents the id of user who want to login.  |
| state | An unguessable random string. It is used to protect against cross-site request forgery attacks. |

> How to do code challenge ?

```go
func main() {
    var codeVerifier = "sampleS256Code"

    codeChallenge := GetCodeChallenge(codeVerifier)
}

func GetCodeChallenge(codeVerifier string) string {
	s256 := sha256.Sum256([]byte(codeVerifier))
	return base64.URLEncoding.EncodeToString(s256[:])
}
```

<br/>

Second, you need to have a GET hook api, that is to say, it’s your `redirect_uri`

You will receive the callback
    
**Callback Body**

| Field | type | Description |
| --- | --- | --- |
| client_id | string| represents your client id. |
| response_type | string | represents the response type (e.g. code) |
| code_challenge | string | the value of code_verifier after sha256 hash. |
| code_challenge_method | string |  |
| user_id | string| represents the id of user who want to login.  |
| state | string|  An unguessable random string. It is used to protect against cross-site request forgery attacks. |
| code | string | the authentication code allow you to do exchange for token |
| redirect_uri | string | represents the callback endpoint |

<br/>

You need to write response body 

1. if you ```success``` to receive callback
    
    **Response Header**
    | Key | Value |
    | --- | --- |
    | Content-Type | application/json |

    **Response Body**
    | Field | Type | Description |
    | --- | --- | --- |
    | success | boolean | ```true``` |
    

2. if you ```fail``` to receive callback

    **Response Header**
    | Key | Value |
    | --- | --- |
    | Content-Type | application/json |

    **Response Body**
    | Field | Type | Description |
    | --- | --- | --- |
    | success | boolean | ```false``` |
    | message | string, optional | error message |

<br />

Then, you can receive authorization code from `redirect_uri`, and use the `code` to exchange the `access_token` by call  OAuth Server `POST /v1/oauth/token`
    
**Header**
| Key | Value |
| --- | --- |
| Content-Type | x-www-form-urlencoded |


**Request Body**
            
| Field | type |Description |
| --- | --- | --- |
| client_id | string| represents your client id. |
| client_secret | string| represents your client secret. |
| code_verifier | string | represents the plain text before doing sha256 hash to code_challenge |
| code | string | the authentication code allow you to do exchange for token |
|grant_type| string | represents the action you want to do, it would be authorization_code on here|


**Response Body**
            
| Field | type |Description |
| --- | --- | --- |
| access_token | string| represents the token you can use to access resources of user |
| token_type | string| represents access_token type (e.g. Bearer) |
| refresh_token | string | represents the token you can use to refresh access_token |
| expires_in | string | represents the time duration of access_token of seconds |

<br/>

### How to refresh access token ?
    
`POST /v1/oauth/token`

**Header**
| Key | Value |
| --- | --- |
| Content-Type | x-www-form-urlencoded |

**Request Body**
        
| Field | type |Description |
| --- | --- | --- |
| client_id | string| represents your client id. |
| client_secret | string| represents your client secret. |
| refresh_token | string | represents the token you can use to refresh access_token |
|grant_type| string | represents the action you want to do, it would be refresh_token on here|



**Response Body**
        
| Field | type |Description |
| --- | --- | --- |
| access_token | string| represents the new token you can use to access resources of user |
| token_type | string| represents access_token type (e.g. Bearer) |
| refresh_token | string | represents the new token you can use to refresh access_token |
| expires_in | string | represents the time duration of access_token of seconds |

<br/>

# Exchange API

**Request Header**

| Key | Value |
| --- | --- |
| Authorization | Bearer [Access Token] | 
| Content-Type | application/json |


<br/>

**Standard Response Format**

| Field | Type | Description |
| --- | --- | --- |
| success | boolean | |
| message | string | error message |
| data | json  | response data |
| request_id | string | represents the request |

<br/>

## Get Account

Query the information of account

Request

```GET /v1/oauth/exchange/account```


Response Format

```json
{
    "success": true,
    "data": {
        "id": "ad122e63-9112-499e-be60-1997f9455f6b",
        "email": "bitgin@bitgin.com",
        "phone": "0912345678",
        "kyc_level": 2
    }
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | ID of account |
| email | string | Email of account |
| phone | string | Phone number of account |
| kyc_level | number | KYC (Know Your Customer) level|

<br/>

## Get Account Bank 

Query the information of bank of account

Request

```GET /v1/oauth/exchange/account/bank```

Response Format

```json
{
    "success": true,
    "data": {
        "bank_code": "004",
        "bank_name": "臺灣銀行",
        "branch_code": "0071",
        "branch_name":"館前分行",
        "holder": "測先生",
        "number": "0060090100005333"
    }
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| bank_code | string | represents code of bank |
| bank_name | string | represents name of bank |
| branch_code | string | represents code of branch |
| branch_name | string | represents name of branch |
| holder | string | represents the holder's name of bank account |
| number | string | represents the number of bank account  |


<br/>


## Get Balance

Query balances of account

Request

```GET /v1/oauth/exchange/wallet/balance?currency={currency}```

Parameters

| Field | Type  | Description |
| :---  | :---  | :---        |
| currency | string | optional, if the field is empty, it will return all balances information as default. (e.g. BTC, ETH, USDT, TWD)|

<br />

Response Format

```json
{
    "success": true,
    "data": [
        {
            "currency": "USDT",
            "total": "45.4",
            "available": "45.4",
            "locked": "0"
        }
    ]
}

```
| Field | Type  | Description |
| :---  | :---  | :---        |
| currency | string | BTC, ETH, USDT, TWD |
| total | decimal | total amount |
| available | decimal | amount available|
| locked | decimal |  amount locked  |

</br>



## Get Deposit History

Query deposit history

Request

```GET /v1/oauth/exchange/wallet/history/deposit?currency={currency}&start_time={start_time}&end_time={end_time}&limit={limit}&offset={offset}```

Parameters 

| Field | Type  | Description |
| :---  | :---  | :---        |
| currency | string | optional, if the field is empty, it will return all deposit history as default. (e.g. BTC, ETH, USDT, TWD)|
| start_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| end_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |
| limit| number | represents limit of pagination | 
| offset | number | represents offset of pagination | 

Response Format

```json
{
    "success": true,
    "data": [
        {
            "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
            "status": "completed",
            "currency": "USDT",
            "chain": "Tron",
            "amount": "1500",
            "fee": "0",
            "fee_currency": "USDT",
            "txid": "ba2f799dd1607a0d118dd9320019ea9ca7e42492760e76abbeb27b29f6404cf7",
            "created_at": 1615974333,
            "updated_at": 1615974333,
            "completed_at": 1615975346
        },
        {
            "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
            "status": "completed",
            "currency": "USDT",
            "chain": "Tron",
            "amount": "750",
            "fee": "0",
            "fee_currency": "USDT",
            "txid": "aa2f799dd1607a0d118dd9320019ea9ca7e42492760e76abbeb27b29f6404ch35",
            "created_at": 1615974333,
            "updated_at": 1615974333,
            "completed_at": 1615975346
        }
    ]
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | represents deposit id|
| [status](#deposit-status-definition) | string | status of deposit|
| currency | string | BTC, ETH, USDT, TWD |
| chain | string | Bitcoin, Ethereum, Tron |
| amount | decimal | total amount |
| fee | decimal |  |
| fee_currency | string | BTC, ETH, USDT, TWD  |
| txid | string | transaction hash |
| created_at | number | when the deposit was created, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| updated_at | number | when the deposit was updated, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| completed_at | number | only exists when the deposit was completed, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |

</br>

## Get Deposit Address

Query deposit addresses of account

Request

```GET /v1/oauth/exchange/wallet/deposit_address?currency={currency}&chain={chain}```

Parameters

| Field | Type  | Description |
| :---  | :---  | :---        |
| currency | string | BTC, ETH, USDT|
| chain | string | optional (e.g. Bitcoin, Ethereum, Tron)|

<br />

Response Format

```json
{
    "success": true,
    "data": [
        {
            "chain": "Tron",
            "address": "TXHzvoDBPaG7YbSgb3zdoosJK4x4Kmf2J2"
        },
        {
            "chain": "Ethereum",
            "address": "0xB76204882Fbef161428588560b48dB570A9d42Bb"
        }
    ]
}

```
| Field | Type  | Description |
| :---  | :---  | :---        |
| chain | string | Bitcoin, Ethereum, Tron |
| address | string |  |

</br>

## Get Deposit Bank 

Query the deposit bank 

Request

```GET /v1/oauth/exchange/wallet/deposit_bank```

Response Format

```json
{
    "success": true,
    "data": {
        "bank_code": "802",
        "bank_name": "凱基商業銀行",
        "branch_code": "0072",
        "branch_name":"城東分行",
        "holder": "凱基商業銀行受託信託財產專戶",
        "number": "51730000007724"
    }
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| bank_code | string | represents code of bank |
| bank_name | string | represents name of bank |
| branch_code | string | represents code of branch |
| branch_name | string | represents name of branch |
| holder | string | represents the holder's name of bank account |
| number | string | represents the number of virtual account  |


<br/>

### Get Withdrawal History

Query withdrawal history

Request

```GET /v1/oauth/exchange/wallet/withdrawal?currency={currency}&start_time={start_time}&end_time={end_time}&limit={limit}&offset={offset}```

Parameters 

| Field | Type  | Description |
| :---  | :---  | :---        |
| currency | string | optional, if the field is empty, it will return all withdrawal history as default. (e.g. BTC, ETH, USDT, TWD)|
| start_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| end_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |
| limit| number | represents limit of pagination | 
| offset | number | represents offset of pagination | 

Response Format

```json
{
    "success": true,
    "data": [
        {
            "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
            "status": "completed",
            "currency": "USDT",
            "chain": "Tron",
            "amount": "1500",
            "fee": "0",
            "fee_currency": "USDT",
            "to_address": "TXHzvoDBPaG7YbSgb3zdoosJK4x4Kmf2J2",
            "txid": "ba2f799dd1607a0d118dd9320019ea9ca7e42492760e76abbeb27b29f6404cf7",
            "created_at": 1615974333,
            "updated_at": 1615974333,
            "completed_at": 1615975346
        },
        {
            "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
            "status": "pending",
            "currency": "USDT",
            "chain": "Tron",
            "amount": "200",
            "fee": "0",
            "fee_currency": "USDT",
            "to_address": "TXHzvoDBPaG7YbSgb3zdoosJK4x4Kmf2J2",
            "created_at": 1615974333,
            "updated_at": 1615974333
        }
    ]
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | represents withdrawal id |
| [status](#withdrawal-status-definition) | string | status of withdrawal|
| currency | string | BTC, ETH, USDT, TWD |
| chain | string | Bitcoin, Ethereum, Tron |
| amount | decimal | total amount |
| fee | decimal |  |
| fee_currency | string | BTC, ETH, USDT, TWD  |
| to_address | string |  |
| txid | string | transaction hash |
| created_at | number | when the withdrawal was created, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| updated_at | number | when the withdrawal was updated, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| completed_at | number | only exists when the withdrawal was completed, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |


</br>

## Request Withdrawal

Send the request of withdrawal

Request

```POST /v1/oauth/exchange/wallet/withdrawal```

Header
| Key | Value |
| --- | --- |
| Content-Type | application/json |

Post Body

```json
{
    "chain": "Tron",
    "currency": "USDT",
    "address": "",
	"amount": "100",
	"is_deduction": true
}
```
| Field | Type  | Description |
| :---  | :---  | :---        |
| chain | string | Bitcoin, Ethereum, Tron |
| currency | string | BTC, ETH, USDT, TWD|
| address | string | withdrawal address |
| amount | decimal |  withdrawal amount |
| is_deduction | boolean |  |

Response Format

```json
{
    "success": true,
    "data": {
        "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
        "status": "pending",
        "currency": "USDT",
        "chain": "Tron",
        "amount": "1500",
        "fee": "0",
        "fee_currency": "USDT",
        "to_address": "TXHzvoDBPaG7YbSgb3zdoosJK4x4Kmf2J2",
        "created_at": 1615974333,
        "updated_at": 1615974333,
    }
}
```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | represents withdrawal id |
| [status](#withdrawal-status-definition) | string | status of withdrawal|
| currency | string | BTC, ETH, USDT, TWD |
| chain | string | Bitcoin, Ethereum, Tron |
| amount | decimal | total amount |
| fee | decimal |  |
| fee_currency | string | BTC, ETH, USDT, TWD  |
| to_address | string | |
| created_at | number | when the withdrawal was created, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| updated_at | number | when the withdrawal was updated, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
<br/>

## Confirm Withdrawal

Confirm the withdrawal with 2FA code

Request 
```POST  /v1/oauth/exchange/wallet/withdrawal/confirm```

Header
| Key | Value |
| --- | --- |
| Content-Type | application/json |

Post Body
```json
{
    "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
    "code": "123456"
}
```

Response Format

```json
{
    "success": true,
    "data": {
        "id": "c82b7de8-c654-4d9e-b84e-022b52c189bb",
        "status": "completed",
        "currency": "USDT",
        "chain": "Tron",
        "amount": "1500",
        "fee": "0",
        "fee_currency": "USDT",
        "type": "crypto",
        "from_address": "TTsNwkygXcdCPxb6BZEkjznGPBDLi5A8pZ",
        "to_address": "TXHzvoDBPaG7YbSgb3zdoosJK4x4Kmf2J2",
        "txid": "ba2f799dd1607a0d118dd9320019ea9ca7e42492760e76abbeb27b29f6404cf7",
        "created_at": 1615974333,
        "updated_at": 1615974333,
        "completed_at": 1615975346,
    }
}
```
| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | represents withdrawal id |
| [status](#withdrawal-status-definition) | string | status of withdrawal|
| currency | string | BTC, ETH, USDT, TWD |
| chain | string | Bitcoin, Ethereum, Tron |
| amount | decimal | total amount |
| fee | decimal |  |
| fee_currency | string | BTC, ETH, USDT, TWD  |
| [type](#withdrawal-type-definition) | string | type of withdrawal|
| from_address | string |  |
| to_address | string |  |
| txid | string | transaction hash |
| created_at | number | when the withdrawal was created, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| updated_at | number | when the withdrawal was updated, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC|
| completed_at | number | only exists when the withdrawal was completed, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |

<br/>

## Get Trade History

Query trade history

Request

```GET /v1/oauth/exchange/trade?market={market}&start_time={start_time}&end_time={end_time}&limit={limit}&offset={offset}```

Parameters 

| Field | Type  | Description |
| :---  | :---  | :---        |
| market | string | optional, if the field is empty, it will return all trade history as default. (e.g. BTCTWD, ETHTWD, USDTTWD)|
| start_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |
| end_time | number | Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |
| limit| number | represents limit of pagination | 
| offset | number | represents offset of pagination | 


Response Format

```json
{
    "success": true,
    "data": [
        {
            "id": "b3422e63-9112-499e-be60-1997f9455f6b",
            "market": "USDTTWD",
            "side": "buy",
            "price": "29.03",
            "size": "1000",
            "fee": "0",
            "fee_rate": "0",
            "fee_currency": "USDT",
            "time": 1615975346
        }
    ]
}

```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | trade id |
| market | string  | BTCTWD, ETHTWD, USDTTWD |
| [side](#order-side-definition) | string  | order side (e.g. buy, sell) |
| price | decimal | price of the transaction|
| size | decimal | total amount of the transaction|
| fee | decimal | fee per transaction|
| fee_rate | decimal | |
| fee_currency | string | Bitcoin, Ethereum, Tron, TWD |
| time | number | when the transaction was completed, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |


</br>


## Get Quote

Get quote of trade

Request

```POST /v1/oauth/exchange/trade/quote```

Header
| Key | Value |
| --- | --- |
| Content-Type | application/json |

Post Body

```json
{
    "market": "USDTTWD",
    "side": "buy",
    "base_amount": "1000",
}
```

| Field | Type  | Description |
| :---  | :---  | :---        |
| market | string | BTCTWD, ETHTWD, USDTTWD|
| [side](#order-side-definition) | string  | order side (e.g. buy, sell) |
| base_amount | decimal | total amount you want to buy/sell (e.g. USDT, ETH, BTC)|
| quote_amount | decimal | total amount you want to buy/sell (e.g. TWD)|

<aside class="notice">
    NOTE: You have to provide base_amount <strong>OR</strong> quote_amount <strong>(only pick one of two)</strong> to ask for quote.
</aside>

Response Format

```json
{
    "success": true,
    "data": {
        "id": "ad122e63-9112-499e-be60-1997f9455f6b",
        "market": "USDTTWD",
        "side": "buy",
        "price": "29.03",
        "base_amount": "1000",
        "quote_amount": "29030",
        "fee": "0",
        "fee_rate": "0",
        "fee_currency": "USDT",
        "created_at": 1615974333,
        "expired_at": 1615974333
    }
}
```

| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | quotation id |
| market | string  | BTCTWD, ETHTWD, USDTTWD |
| [side](#order-side-definition) | string  | order side (e.g. buy, sell) |
| price | decimal | current price |
| base_amount | decimal | total amount you want to buy/sell (e.g. USDT, ETH, BTC) |
| quote_amount | decimal | total amount you want to buy/sell (e.g. TWD) |
| fee | decimal | fee per transaction|
| fee_rate | decimal | |
| fee_currency | string | Bitcoin, Ethereum, Tron, TWD |
| created_at | number | when the quotation was created, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |
| expired_at | number | when the quotation expires, Unix time of current time, the number of `milliseconds` elapsed since January 1, 1970 UTC |

<aside class="notice">
    NOTE: For complete the transaction, you can use <strong>Accept Quote API</strong> to transaction by using the quotation id.
</aside>

</br>


### Accept Quote

Accept quote of trade

Request

```POST /v1/oauth/exchange/trade/quote/accept```

Header
| Key | Value |
| --- | --- |
| Content-Type | application/json |

Post Body

```json
{
    "id": "ad122e63-9112-499e-be60-1997f9455f6b"
}
```

<aside class="notice">
    NOTE: You can complete trade with quote id which is from <strong>Get Quote API
    </strong> by using <strong>Accept Quote API</strong>
</aside>


| Field | Type  | Description |
| :---  | :---  | :---        |
| id | string | quote id |

Response Format

```json
{
    "success": true
}
```
### Withdrawal Type Definition
| Value | Description |
| :---  | :---     |
| crypto  | crypto withdrawal|
| fiat_kgi  | fiat withdrawal |
| internal_transfer  | internal transfer (e.g. from BITGIN address to BITGIN address)|


## Withdrawal Status Definition

| Value | Description |
| :---  | :---     |
| pending  | withdrawal is waiting to be sent |
| waiting_approval  | withdrawal is waiting for an approval |
| approved  | withdrawal approved|
| bank_verifying  | withdrawal is in bank procedure |
| sent  |  withdrawal sent |
| completed  | withdrawal completed |
| cancelled  | withdrawal has been cancelled |
| rejected  |  withdrawal has been rejected |
| failed  | withdrawal failed |

## Deposit Status Definition

| Value | Description |
| :---  | :---     |
| waiting_confirmation	 | deposit confirmation count on the blockchain |
| completed  | deposit completed |
| rejected  | deposit has been rejected |

## Order Side Definition

| Value | Description |
| :---  | :---     |
| buy  | order side buy |
| sell  | order side sell |