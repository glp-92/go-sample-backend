# API Endpoints

## Auth

### Login

[POST] `{url}/login`

Body:
- username: str
- password: str,
- topt: str

### Logout

[POST] `{url}/revoke`

Authorization: Bearer token

### Validate token

[GET] `{url}/valid`

## Categories

### Get all categories

[GET] `{url}/categories` - should change `category`

### Get categories pageable

[GET] `{url}/categories?page={page}`
`&perpage={perpage}` - should change?

### Create category

[POST] `{url}/categories`

Authorization: Bearer token

Body:
- name: str
- slug: str

### Delete category

[DELETE] `{url}/categories/{categoryId}`

Authorization: Bearer token

### Update category

[PUT] `{url}/categories`

Authorization: Bearer token

Body: category (revisar)

