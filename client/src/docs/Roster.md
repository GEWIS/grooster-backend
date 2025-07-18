# Roster


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdAt** | **string** |  | [optional] [default to undefined]
**date** | **string** |  | [optional] [default to undefined]
**deletedAt** | [**GormDeletedAt**](GormDeletedAt.md) |  | [optional] [default to undefined]
**id** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**organ** | [**Organ**](Organ.md) |  | [optional] [default to undefined]
**organId** | **number** |  | [optional] [default to undefined]
**rosterAnswer** | [**Array&lt;RosterAnswer&gt;**](RosterAnswer.md) |  | [optional] [default to undefined]
**rosterShift** | [**Array&lt;RosterShift&gt;**](RosterShift.md) |  | [optional] [default to undefined]
**saved** | **boolean** |  | [optional] [default to undefined]
**updatedAt** | **string** |  | [optional] [default to undefined]
**values** | **Array&lt;string&gt;** |  | [optional] [default to undefined]

## Example

```typescript
import { Roster } from './api';

const instance: Roster = {
    createdAt,
    date,
    deletedAt,
    id,
    name,
    organ,
    organId,
    rosterAnswer,
    rosterShift,
    saved,
    updatedAt,
    values,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
