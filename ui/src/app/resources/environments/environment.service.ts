import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { Environment } from './environment.component';
import { AuthService } from '../../auth/auth.service';

@Injectable()
export class EnvironmentService {

  private environmentsCreateUrl = '/api/v1/environments'
  private environmentsGetUrl    = '/api/v1/environments/'
  private environmentsEditUrl   = '/api/v1/environments/'

  constructor(private http: Http, private auth: AuthService) { }

  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }

  createEnvironment(name: string, applicationId: number): Observable<Environment> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, applicationId: applicationId }

    return this.http.post(this.environmentsCreateUrl, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

    getEnvironment(uuid: string): Observable<Environment> {
      return this.http.get(this.environmentsGetUrl + uuid, this.auth.jwtRequestOptions())
        .map((res: Response) => res.json())
        .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
    }

    updateEnvironment(environment: Environment): Observable<Environment> {
      const options = new RequestOptions({ headers: this.getHeaders() });
      const data = { name: environment.name };

      return this.http.patch(this.environmentsEditUrl + environment.uuid, data, this.auth.jwtRequestOptions())
        .map((res: Response) => res.json())
        .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
    }
}
