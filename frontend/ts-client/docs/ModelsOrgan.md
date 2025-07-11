# ModelsOrgan

An organ that users can be part of.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdAt** | **string** |  | [optional] [default to undefined]
**deletedAt** | [**GormDeletedAt**](GormDeletedAt.md) |  | [optional] [default to undefined]
**id** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**updatedAt** | **string** |  | [optional] [default to undefined]
**users** | [**Array&lt;User&gt;**](User.md) |  | [optional] [default to undefined]

## Example

```typescript
import { ModelsOrgan } from './api';

const instance: ModelsOrgan = {
    createdAt,
    deletedAt,
    id,
    name,
    updatedAt,
    users,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
