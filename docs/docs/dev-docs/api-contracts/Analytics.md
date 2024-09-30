# Admin Analytics

## Helper types

```ts
interface time {
    days : number;
    hours : number;
    minutes: number;
}
```
```json
 //idk what the ts of this is
 // This is a GroupedValue format
 {
    "key" : number,
    ...
 }
```

# Predicted time to complete a dispute

- **Endpoint:** `GET /time/estimation`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** only available to admins

```ts 
//returns
//if resolved disputes exist
interface ResponseResolvedDisputes {
    data : time
}
//otherwise
interface error {
    error : string
}
```

# Count disputes by country

- **Endpoint:** `GET /dispute/countries`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** only available to admins

```ts 
//returns
interface CountriesDisputeCount {
    data: GroupedValue
} 
```

# Aggregate analytics

- **Endpoint:** `POST /stats/{table}`
- **Headers:**
  - `Authorization: Bearer <JWT>`
- **Note:** only available to admins

```ts 
//queries the count of records, applying the below params if supplied
interface request{
    group? : string;
    where? : {
        column : string,
        value : string
    }
}

//returns
interface CountriesDisputeCount {
    data: GroupedValue
} 
```



