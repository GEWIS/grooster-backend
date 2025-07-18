# RosterApi

All URIs are relative to _http://localhost_

| Method                                            | HTTP request                     | Description                                                       |
| ------------------------------------------------- | -------------------------------- | ----------------------------------------------------------------- |
| [**createRoster**](#createroster)                 | **POST** /roster                 | CreateRoster a new roster                                         |
| [**createRosterTemplate**](#createrostertemplate) | **POST** /roster/template        | Creates a template of a roster by defining the name of the shifts |
| [**deleteRoster**](#deleteroster)                 | **DELETE** /roster/{id}          | DeleteRoster a roster                                             |
| [**deleteRosterTemplate**](#deleterostertemplate) | **DELETE** /roster/template/{id} | Deletes a roster template by ID                                   |
| [**getRoster**](#getroster)                       | **GET** /roster/{id}             | Get a specific roster by id                                       |
| [**getRosterTemplate**](#getrostertemplate)       | **GET** /roster/template/{id}    | Get a roster template by ID                                       |
| [**getRosterTemplates**](#getrostertemplates)     | **GET** /roster/template         | Get all rosters templates or query by organ ID                    |
| [**getRosters**](#getrosters)                     | **GET** /roster                  | Get all rosters or query by date and organ                        |
| [**updateRoster**](#updateroster)                 | **PATCH** /roster/{id}           | Update a roster                                                   |

# **createRoster**

> Roster createRoster(createParams)

### Example

```typescript
import { RosterApi, Configuration, RosterCreateRequest } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let createParams: RosterCreateRequest; //Roster input

const { status, data } = await apiInstance.createRoster(createParams);
```

### Parameters

| Name             | Type                    | Description  | Notes |
| ---------------- | ----------------------- | ------------ | ----- |
| **createParams** | **RosterCreateRequest** | Roster input |       |

### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRosterTemplate**

> Array<SavedShift> createRosterTemplate()

### Example

```typescript
import { RosterApi, Configuration, RosterTemplateCreateRequest } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let params: RosterTemplateCreateRequest; //Template Params (optional)

const { status, data } = await apiInstance.createRosterTemplate(params);
```

### Parameters

| Name       | Type                            | Description     | Notes |
| ---------- | ------------------------------- | --------------- | ----- |
| **params** | **RosterTemplateCreateRequest** | Template Params |       |

### Return type

**Array<SavedShift>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details

| Status code | Description      | Response headers |
| ----------- | ---------------- | ---------------- |
| **200**     | Created Template | -                |
| **400**     | Invalid request  | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRoster**

> string deleteRoster()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.deleteRoster(id);
```

### Parameters

| Name   | Type         | Description | Notes                 |
| ------ | ------------ | ----------- | --------------------- |
| **id** | [**number**] | Roster ID   | defaults to undefined |

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |
| **404**     | Not Found   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRosterTemplate**

> string deleteRosterTemplate()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Template ID (default to undefined)

const { status, data } = await apiInstance.deleteRosterTemplate(id);
```

### Parameters

| Name   | Type         | Description | Notes                 |
| ------ | ------------ | ----------- | --------------------- |
| **id** | [**number**] | Template ID | defaults to undefined |

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |
| **404**     | Not Found   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRoster**

> Roster getRoster()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.getRoster(id);
```

### Parameters

| Name   | Type         | Description | Notes                 |
| ------ | ------------ | ----------- | --------------------- |
| **id** | [**number**] | Roster ID   | defaults to undefined |

### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |
| **404**     | Not Found   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterTemplate**

> RosterTemplate getRosterTemplate()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Template ID (default to undefined)

const { status, data } = await apiInstance.getRosterTemplate(id);
```

### Parameters

| Name   | Type         | Description | Notes                 |
| ------ | ------------ | ----------- | --------------------- |
| **id** | [**number**] | Template ID | defaults to undefined |

### Return type

**RosterTemplate**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |
| **404**     | Not Found   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterTemplates**

> Array<RosterTemplate> getRosterTemplates()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let organId: number; // (optional) (default to undefined)

const { status, data } = await apiInstance.getRosterTemplates(organId);
```

### Parameters

| Name        | Type         | Description | Notes                            |
| ----------- | ------------ | ----------- | -------------------------------- |
| **organId** | [**number**] |             | (optional) defaults to undefined |

### Return type

**Array<RosterTemplate>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosters**

> Array<Roster> getRosters()

### Example

```typescript
import { RosterApi, Configuration } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let date: string; //Date filter (ISO format) (optional) (default to undefined)
let organId: number; //Organ ID filter (optional) (default to undefined)

const { status, data } = await apiInstance.getRosters(date, organId);
```

### Parameters

| Name        | Type         | Description              | Notes                            |
| ----------- | ------------ | ------------------------ | -------------------------------- |
| **date**    | [**string**] | Date filter (ISO format) | (optional) defaults to undefined |
| **organId** | [**number**] | Organ ID filter          | (optional) defaults to undefined |

### Return type

**Array<Roster>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRoster**

> Roster updateRoster(updateParams)

### Example

```typescript
import { RosterApi, Configuration, RosterUpdateRequest } from "./api";

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)
let updateParams: RosterUpdateRequest; //Roster input

const { status, data } = await apiInstance.updateRoster(id, updateParams);
```

### Parameters

| Name             | Type                    | Description  | Notes                 |
| ---------------- | ----------------------- | ------------ | --------------------- |
| **updateParams** | **RosterUpdateRequest** | Roster input |                       |
| **id**           | [**number**]            | Roster ID    | defaults to undefined |

### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |
| **400**     | Bad Request | -                |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)
