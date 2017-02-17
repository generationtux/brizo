import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { AuthService } from '../auth/auth.service';
import { Version } from './version.model';

@Injectable()
export class VersionService {

  private url = '/api/v1/environments/';

  constructor(private http: Http, private auth: AuthService) {}

  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }

  createVersion(version: Version): Observable<Version> {
    let options = new RequestOptions({ headers: this.getHeaders() });
    let url = this.url + version.environment_uuid + '/versions';
    let data = {
      name: version.name,
      image: version.image,
      replicas: version.replicas,
      volumes: version.volumes,
      pullPolicy: version.pullPolicy,
      args: version.args,
      ports: version.ports,
      volumeMounts: version.volumeMounts,
    }

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

  getVersion(environmentUuid: string, uuid: string): Observable<Version> {
    const url = this.url + environmentUuid + '/versions/' + uuid;

    return this.http.get(url, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  updateVersion(environmentUuid: string, uuid: string, name: string): Observable<Version> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name };

    const url = this.url + environmentUuid + '/versions/' + uuid;

    return this.http.patch(url, data, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

}
