import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { AuthService } from '../../auth/auth.service';

@Injectable()
export class ConfigurationService {

  private url = '/api/v1/environments/';

  constructor(private http: Http, private auth: AuthService) { }

  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }

  createConfiguration(name: string, value: string, environmentId: number, environmentUuid: string): Observable<any> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, value: value, environment_id: environmentId }
    
    const url = this.url + environmentUuid + '/configuration';

    return this.http.post(url, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}