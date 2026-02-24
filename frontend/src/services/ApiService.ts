import {
  RosterApi,
  UserApi,
  AuthApi,
  RosterAnswerApi,
  RosterShiftApi,
  SavedShiftApi,
  Configuration,
  ExportApi,
  OrganApi,
  ShiftGroupApi,
} from '@gewis/grooster-backend-ts';
import { getTokenFromStorage } from '@/helpers/TokenHelper';

class ApiService {
  private readonly _rosterApi: RosterApi;

  private readonly _rosterAnswerApi: RosterAnswerApi;

  private readonly _rosterShiftApi: RosterShiftApi;

  private readonly _savedShiftApi: SavedShiftApi;

  private readonly _shiftGroupApi: ShiftGroupApi;

  private readonly _exportApi: ExportApi;

  private readonly _userApi: UserApi;

  private readonly _organApi: OrganApi;

  private readonly _authApi: AuthApi;

  constructor(basePath: string) {
    const config = new Configuration({
      basePath,
      apiKey: () => {
        const token = getTokenFromStorage().token;
        return token ? `Bearer ${token}` : '';
      },
    });

    this._rosterApi = new RosterApi(config);
    this._rosterAnswerApi = new RosterAnswerApi(config);
    this._rosterShiftApi = new RosterShiftApi(config);
    this._savedShiftApi = new SavedShiftApi(config);
    this._shiftGroupApi = new ShiftGroupApi(config);
    this._exportApi = new ExportApi(config);
    this._userApi = new UserApi(config);
    this._organApi = new OrganApi(config);
    this._authApi = new AuthApi(config);
  }

  get roster(): RosterApi {
    return this._rosterApi;
  }

  get rosterAnswer(): RosterAnswerApi {
    return this._rosterAnswerApi;
  }

  get rosterShift(): RosterShiftApi {
    return this._rosterShiftApi;
  }

  get savedShiftApi(): SavedShiftApi {
    return this._savedShiftApi;
  }

  get shiftGroupApi(): ShiftGroupApi {
    return this._shiftGroupApi;
  }

  get exportApi(): ExportApi {
    return this._exportApi;
  }

  get user(): UserApi {
    return this._userApi;
  }

  get organ(): OrganApi {
    return this._organApi;
  }

  get auth(): AuthApi {
    return this._authApi;
  }
}

const isGrooster = /grooster(\.\w+)?\.gewis\.nl$/.test(window.location.hostname);

const location = isGrooster ? window.location.origin : 'http://localhost:8080';

const apiService = new ApiService(location + '/api/v1');

export default apiService;
