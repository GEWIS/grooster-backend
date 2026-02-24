# RosterApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createRoster**](#createroster) | **POST** /roster | CreateRoster a new roster|
|[**createRosterTemplate**](#createrostertemplate) | **POST** /roster/template | Creates a template of a roster by defining the name of the shifts|
|[**createRosterTemplateShiftPreference**](#createrostertemplateshiftpreference) | **POST** /roster/template/shift-preference | Creates a roster template shift preference|
|[**deleteRoster**](#deleteroster) | **DELETE** /roster/{id} | DeleteRoster a roster|
|[**deleteRosterTemplate**](#deleterostertemplate) | **DELETE** /roster/template/{id} | Deletes a roster template by ID|
|[**getRoster**](#getroster) | **GET** /roster/{id} | Get a specific roster by id|
|[**getRosterTemplate**](#getrostertemplate) | **GET** /roster/template/{id} | Get a roster template by ID|
|[**getRosterTemplateShiftPreferences**](#getrostertemplateshiftpreferences) | **GET** /roster/template/shift-preference | Gets shift preferences filtered by user and template|
|[**getRosterTemplates**](#getrostertemplates) | **GET** /roster/template | Get all rosters templates or query by organ ID|
|[**getRosters**](#getrosters) | **GET** /roster | Get all rosters or query by date and organ|
|[**updateRoster**](#updateroster) | **PATCH** /roster/{id} | Update a roster|
|[**updateRosterTemplate**](#updaterostertemplate) | **PUT** /roster/template/{id} | Updates a roster template by ID|
|[**updateRosterTemplateShift**](#updaterostertemplateshift) | **PATCH** /roster/template/shift/{id} | Updates a roster template shift by ID|
|[**updateRosterTemplateShiftPreference**](#updaterostertemplateshiftpreference) | **PATCH** /roster/template/shift-preference/{id} | Updates a roster template shift preference by ID|

# **createRoster**
> Roster createRoster(createParams)


### Example

```typescript
import {
    RosterApi,
    Configuration,
    RosterCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let createParams: RosterCreateRequest; //Roster input

const { status, data } = await apiInstance.createRoster(
    createParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createParams** | **RosterCreateRequest**| Roster input | |


### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRosterTemplate**
> Array<RosterTemplate> createRosterTemplate()


### Example

```typescript
import {
    RosterApi,
    Configuration,
    TemplateCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let params: TemplateCreateRequest; //Template Params (optional)

const { status, data } = await apiInstance.createRosterTemplate(
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **TemplateCreateRequest**| Template Params | |


### Return type

**Array<RosterTemplate>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Created Template |  -  |
|**400** | Invalid request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRosterTemplateShiftPreference**
> RosterTemplateShiftPreference createRosterTemplateShiftPreference(params)


### Example

```typescript
import {
    RosterApi,
    Configuration,
    TemplateShiftPreferenceCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let params: TemplateShiftPreferenceCreateRequest; //Creation params

const { status, data } = await apiInstance.createRosterTemplateShiftPreference(
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **TemplateShiftPreferenceCreateRequest**| Creation params | |


### Return type

**RosterTemplateShiftPreference**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRoster**
> string deleteRoster()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.deleteRoster(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster ID | defaults to undefined|


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRosterTemplate**
> string deleteRosterTemplate()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Template ID (default to undefined)

const { status, data } = await apiInstance.deleteRosterTemplate(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Template ID | defaults to undefined|


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRoster**
> Roster getRoster()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)

const { status, data } = await apiInstance.getRoster(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Roster ID | defaults to undefined|


### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterTemplate**
> RosterTemplate getRosterTemplate()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Template ID (default to undefined)

const { status, data } = await apiInstance.getRosterTemplate(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Template ID | defaults to undefined|


### Return type

**RosterTemplate**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterTemplateShiftPreferences**
> Array<RosterTemplateShiftPreference> getRosterTemplateShiftPreferences()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let userId: number; //User ID (default to undefined)
let templateId: number; //Template ID (default to undefined)

const { status, data } = await apiInstance.getRosterTemplateShiftPreferences(
    userId,
    templateId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**number**] | User ID | defaults to undefined|
| **templateId** | [**number**] | Template ID | defaults to undefined|


### Return type

**Array<RosterTemplateShiftPreference>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosterTemplates**
> Array<RosterTemplate> getRosterTemplates()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let organId: number; // (optional) (default to undefined)

const { status, data } = await apiInstance.getRosterTemplates(
    organId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **organId** | [**number**] |  | (optional) defaults to undefined|


### Return type

**Array<RosterTemplate>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRosters**
> Array<Roster> getRosters()


### Example

```typescript
import {
    RosterApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let date: string; // (optional) (default to undefined)
let id: number; // (optional) (default to undefined)
let organId: number; // (optional) (default to undefined)

const { status, data } = await apiInstance.getRosters(
    date,
    id,
    organId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **date** | [**string**] |  | (optional) defaults to undefined|
| **id** | [**number**] |  | (optional) defaults to undefined|
| **organId** | [**number**] |  | (optional) defaults to undefined|


### Return type

**Array<Roster>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRoster**
> Roster updateRoster(updateParams)


### Example

```typescript
import {
    RosterApi,
    Configuration,
    RosterUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Roster ID (default to undefined)
let updateParams: RosterUpdateRequest; //Roster input

const { status, data } = await apiInstance.updateRoster(
    id,
    updateParams
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateParams** | **RosterUpdateRequest**| Roster input | |
| **id** | [**number**] | Roster ID | defaults to undefined|


### Return type

**Roster**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRosterTemplate**
> RosterTemplate updateRosterTemplate()


### Example

```typescript
import {
    RosterApi,
    Configuration,
    TemplateUpdateParams
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Template ID (default to undefined)
let params: TemplateUpdateParams; //Update params (optional)

const { status, data } = await apiInstance.updateRosterTemplate(
    id,
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **TemplateUpdateParams**| Update params | |
| **id** | [**number**] | Template ID | defaults to undefined|


### Return type

**RosterTemplate**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRosterTemplateShift**
> RosterTemplateShift updateRosterTemplateShift(params)


### Example

```typescript
import {
    RosterApi,
    Configuration,
    TemplateShiftUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Shift ID (default to undefined)
let params: TemplateShiftUpdateRequest; //Update params

const { status, data } = await apiInstance.updateRosterTemplateShift(
    id,
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **TemplateShiftUpdateRequest**| Update params | |
| **id** | [**number**] | Shift ID | defaults to undefined|


### Return type

**RosterTemplateShift**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRosterTemplateShiftPreference**
> RosterTemplateShiftPreference updateRosterTemplateShiftPreference(params)


### Example

```typescript
import {
    RosterApi,
    Configuration,
    TemplateShiftPreferenceUpdateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new RosterApi(configuration);

let id: number; //Preference ID (default to undefined)
let params: TemplateShiftPreferenceUpdateRequest; //Update params

const { status, data } = await apiInstance.updateRosterTemplateShiftPreference(
    id,
    params
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **params** | **TemplateShiftPreferenceUpdateRequest**| Update params | |
| **id** | [**number**] | Preference ID | defaults to undefined|


### Return type

**RosterTemplateShiftPreference**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

