import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { Environment } from './environment.component';
import { AuthService } from '../../auth/auth.service';

@Injectable()
export class EnvironmentService {

  private environmentsCreateUrl = '/api/v1/environments'

  constructor(private http: Http, private auth: AuthService) { }

  createEnvironment(name: string, applicationId: number): Observable<Environment> {
    const headers = new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
    const options = new RequestOptions({ headers: headers });
    const data = { name: name, applicationId: applicationId }

    return this.http.post(this.environmentsCreateUrl, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}
