## Endpoints

Available endpoints are listed here

### Register

**POST** - `/api/v1/links`

**Body:**

| Name        | Type   | Default              | Required |
|-------------|--------|----------------------|----------|
| `username`  | string | `"test"`             | Yes      |
| `password`  | string | `"somesafepassword"` | Yes      |
| `email`     | string | `"some@email.com"`   | No       |
| `firstName` | string | Iman                 | No       |
| `lastName`  | string | Daneshi              | No       |

**Headers:**

| Name              | Type | Default     | Required |
|-------------------|------|-------------|----------|
| `Authorization` | str  | `"Token somevalidtoken"`  | Yes      |