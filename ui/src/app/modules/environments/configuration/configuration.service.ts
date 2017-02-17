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
  
  getConfiguration(environmentUuid: string): Observable<any> {
    const url = this.url + environmentUuid + '/configuration';
    
    return this.http.get(url, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  createConfiguration(name: string, value: string, environmentUuid: string): Observable<any> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, value: value, environmentUUID: environmentUuid }
    
    const url = this.url + environmentUuid + '/configuration';

    return this.http.post(url, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  deleteConfiguration(environmentId: number, environmentUuid: string): Observable<any> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const url = this.url + environmentUuid + '/configuration/' + environmentId;
    
    return this.http.delete(url, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}