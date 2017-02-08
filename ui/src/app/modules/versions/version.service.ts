import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { AuthService } from '../auth/auth.service';
import { Version } from './version.model';

@Injectable()
export class VersionService {
  
  private url = '/api/v1/environments/';
  //private versionsCreateUrl = '/api/v1/versions'
  //private versionsGetUrl    = '/api/v1/versions/'
  //private versionsEditUrl   = '/api/v1/versions/'
  //private versionsIndexUrl  = '/api/v1/environments/:environment-uuid/versions'
  
  constructor(private http: Http, private auth: AuthService) {}
  
  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }
  
  createVersion(name: string, image: string, replicas: number, environmentId: number, environmentUuid: string): Observable<Version> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, image: image, replicas: replicas, environment_id: environmentId }
    
    const url = this.url + environmentUuid + '/versions';

    return this.http.post(url, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  getVersions(environmentUuid: string): Observable<Version[]> {
    const url = this.url + environmentUuid + '/versions';
    
    return this.http.get(url, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  getVersion(uuid: string): Observable<Version> {
    return this.http.get(this.url + uuid, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  updateVersion(version: Version): Observable<Version> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: version.name };

    return this.http.patch(this.url + version.uuid, data, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
}