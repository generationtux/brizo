import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { AuthService } from '../auth/auth.service';
import { Version } from './version.model';

@Injectable()
export class VersionService {
  
  private versionsCreateUrl = '/api/v1/versions'
  private versionsGetUrl    = '/api/v1/versions/'
  private versionsEditUrl   = '/api/v1/versions/'
  
  constructor(private http: Http, private auth: AuthService) {}
  
  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }
  
  createVersion(name: string, applicationId: number): Observable<Version> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, applicationId: applicationId }

    return this.http.post(this.versionsCreateUrl, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  getVersion(uuid: string): Observable<Environment> {
    return this.http.get(this.versionsGetUrl + uuid, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
  updateVersion(version: Version): Observable<Version> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: version.name };

    return this.http.patch(this.versionsEditUrl + version.uuid, data, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  
}