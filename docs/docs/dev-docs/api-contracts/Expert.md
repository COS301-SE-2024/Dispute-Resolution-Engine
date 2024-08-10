# Experts

# Get Expert

- **Endpoint:** `GET /expert`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface GetExperts {}
```

# suggest Assignment

- **Endpoint:** `GET /expert/suggest-assign`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface SuggestAssign {}
```

# Expert Assignment

- **Endpoint:** `GET /expert/assign`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface ExpertAssign {}
```

# ReAssign Expert

- **Endpoint:** `GET /expert/reassign`
- **Headers:**
  - `Authorization: Bearer <JWT>`

```ts
interface ReAssignExperts {}
```
