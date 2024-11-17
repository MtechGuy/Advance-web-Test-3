# API Usage Guide
<details>
<summary> ## REGISTER USER </summary>

### Step 1: Register a New User
```bash
# Define the request body
BODY='{"username": "John Doe", "email": "john@example.com", "password": "mangotree"}'

# Register the user
curl -d "$BODY" http://localhost:4000/api/v1/register/user
```

### Step 2: Activate the User **NOTE: the TOKEN VALUE is sent through email**
```bash
# Replace "TOKEN_VALUE" with the token sent via email
curl -X PUT -d '{"token": "TOKEN_VALUE"}' http://localhost:4000/api/v1/users/activated
```
</details>
<details>

<summary> ## AUTHENTICATE THE USER </summary>

### Step 1: Request authentication by logging in to authentication endpoint

```bash
# Define the request body
BODY='{"email": "john@example.com", "password": "mangotree"}'

# Trigger the authentication endpoint
curl -d "$BODY" http://localhost:4000/api/v1/tokens/authentication
```

### Step 2: Use the "BEARER TOKEN"
 ```bash
# Replace "BEARER_TOKEN" with the token returned in the previous step
curl -i -H "Authorization: Bearer BEARER_TOKEN" http://localhost:4000/api/v1/healthcheck
```

### Step 3: fetch a specific user Information **WITHOUT PASSWORD**
``` bash
# Replace ":uid" with the user ID
curl -i -H "Authorization: Bearer BEARER_TOKEN" http://localhost:4000/api/v1/users/uid
```
</details>
<details>

</summary>