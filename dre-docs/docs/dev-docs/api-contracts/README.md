# API Contracts

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
  error: "error message",
}
```
