All endpoints follow the following general type:
```ts
type Result<T> =
  | {
      data: T;
      error?: never;
    }
  | {
      data?: never;
      error: string;
    };
```

Which corresponds to either returning:
```json5
{
    "data": /* ... some data */
}
```
or
```json5
{
    "error": "error message"
}
```
# Get Expert
- **Endpoint:** `GET /expert`
- **Headers:** 
  - `Authorization: Bearer <JWT>`

```ts
interface GetExperts{

}
```

# suggest Assignment
- **Endpoint:** `GET /expert/suggest-assign`
- **Headers:** 
  - `Authorization: Bearer <JWT>`

```ts
interface SuggestAssign{

}
```

# Expert Assignment
- **Endpoint:** `GET /expert/assign`
- **Headers:** 
  - `Authorization: Bearer <JWT>`

```ts
interface ExpertAssign{

}
```

# ReAssign Expert
- **Endpoint:** `GET /expert/reassign`
- **Headers:** 
  - `Authorization: Bearer <JWT>`

```ts
interface ReAssignExperts{
    
}
```