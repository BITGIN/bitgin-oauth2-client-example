# OAuth2 Client

**Standard Response**

| Field | Type | Description |
| --- | --- | --- |
| success | boolean | |
| message | string | error message |
| data | json  | response data |
| request_id | string | represents the request |

<br/>

### How to get Access Token ?

First of all, the entry point of OAuth2 Authorization process, images you have a button on your application UI, when user trigger the button, 

- it will be redirect to  BITGIN Domain `GET /oauth/authorize`  with the following parameters
    - Parameters
        
        | Field | Description |
        | --- | --- |
        | client_id | represents your client id. |
        | code_challenge | the value of code_verifier after sha256 hash. |
        | user_id | represents the id of user who want to login.  |
        | state | An unguessable random string. It is used to protect against cross-site request forgery attacks. |

<br/>

Second, you need to have a GET hook api, that is to say, it’s your `redirect_uri`

- You will receive the callback
    
    - Callback Body

        | Field | type |Description |
        | --- | --- | --- |
        | client_id | string| represents your client id. |
        | response_type | string | represents the response type (e.g. code) |
        | code_challenge | string | the value of code_verifier after sha256 hash. |
        | code_challenge_method | string |  |
        | user_id | string| represents the id of user who want to login.  |
        | state | string|  An unguessable random string. It is used to protect against cross-site request forgery attacks. |
        | code | string | the authentication code allow you to do exchange for token |
        | redirect_uri | string | represents the callback endpoint |
        
    Then, you can receive authorization code from `redirect_uri`, and use the `code` to exchange the `access_token` by 
    
    call  OAuth Server `POST /v1/oauth/token`
    
    - Content-Type: x-www-form-urlencoded
    - Request Body
        
        | Field | type |Description |
        | --- | --- | --- |
        | client_id | string| represents your client id. |
        | client_secret | string| represents your client secret. |
        | code_verifier | string | represents the plain text before doing sha256 hash to code_challenge |
        | code | string | the authentication code allow you to do exchange for token |
        |grant_type| string | represents the action you want to do, it would be authorization_code on here|
  
        
    - Response Body
        
        | Field | type |Description |
        | --- | --- | --- |
        | access_token | string| represents the token you can use to access resources of user |
        | token_type | string| represents access_token type (e.g. Bearer) |
        | refresh_token | string | represents the token you can use to refresh access_token |
        | expires_in | string | represents the time duration of access_token of seconds |

        
- A the end, You need to write response body `OK` string if you get the access token successfully.

<br />

### How to Refresh Access Token ?
    
- `POST /v1/oauth/token`
    - Content-Type: x-www-form-urlencoded
    - Request Body
        
        | Field | type |Description |
        | --- | --- | --- |
        | client_id | string| represents your client id. |
        | client_secret | string| represents your client secret. |
        | refresh_token | string | represents the token you can use to refresh access_token |
        |grant_type| string | represents the action you want to do, it would be refresh_token on here|



    - Response Body
        
        | Field | type |Description |
        | --- | --- | --- |
        | access_token | string| represents the new token you can use to access resources of user |
        | token_type | string| represents access_token type (e.g. Bearer) |
        | refresh_token | string | represents the new token you can use to refresh access_token |
        | expires_in | string | represents the time duration of access_token of seconds |

- How to do code challenge ?
    - golang example
        ```go
            /*
            * SampleCodeChallenge
            */
            func SampleCodeChallenge(codeChallenge, codeVerifier string) bool {
                s256 := sha256.Sum256([]byte(codeVerifier))
                // trim padding
                a := strings.TrimRight(base64.URLEncoding.EncodeToString(s256[:]), "=")
                b := strings.TrimRight(codeChallenge, "=")

                return a == b
            }
        ```