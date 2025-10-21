import {
  RosterApi,
  UserApi,
  AuthApi,
  RosterAnswerApi,
  RosterShiftApi,
  SavedShiftApi,
 Configuration } from '@gewis/grooster-backend-ts';
import { getTokenFromStorage } from '@/helpers/TokenHelper';

class ApiService {
  private readonly _rosterApi: RosterApi;

  private readonly _rosterAnswerApi: RosterAnswerApi;

  private readonly _rosterShiftApi: RosterShiftApi;

  private readonly _savedShiftApi: SavedShiftApi;

  private readonly _userApi: UserApi;

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
    this._userApi = new UserApi(config);
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

  get user(): UserApi {
    return this._userApi;
  }

  get auth(): AuthApi {
    return this._authApi;
  }
}

const apiService = new ApiService(window.location.origin + '/api/v1');

export default apiService;
