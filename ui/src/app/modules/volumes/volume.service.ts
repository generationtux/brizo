import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { AuthService } from '../auth/auth.service';
import { Volume } from './volume.model';

@Injectable()
export class VolumeService {

  constructor(private http: Http, private auth: AuthService) {}

  private buildResourceUrl(environmentUuid: string, versionUuid: string): string {
    return '/api/v1/environments/' + environmentUuid + '/versions/' + versionUuid + '/volumes'
  }

  private getHeaders() {
    return new Headers({
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + this.auth.getToken(),
    });
  }

  createVolume(name: string, type: string, source: string, environmentUuid: string, versionUuid: string): Observable<Volume> {
    const options = new RequestOptions({ headers: this.getHeaders() });
    const data = { name: name, type: type, source: source };
    const url = this.buildResourceUrl(environmentUuid, versionUuid);

    return this.http.post(url, data, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  getVolumes(environmentUuid: string, versionUuid: string): Observable<Volume[]> {
    const url = this.buildResourceUrl(environmentUuid, versionUuid);

    return this.http.get(url, this.auth.jwtRequestOptions())
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
  //
  // getVersion(environmentUuid: string, uuid: string): Observable<Version> {
  //   const url = this.url + environmentUuid + '/versions/' + uuid;
  //
  //   return this.http.get(url, this.auth.jwtRequestOptions())
  //     .map((res: Response) => res.json())
  //     .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  // }
  //
  // updateVersion(environmentUuid: string, uuid: string, name: string): Observable<Version> {
  //   const options = new RequestOptions({ headers: this.getHeaders() });
  //   const data = { name: name };
  //
  //   const url = this.url + environmentUuid + '/versions/' + uuid;
  //
  //   return this.http.patch(url, data, this.auth.jwtRequestOptions())
  //     .map((res: Response) => res.json())
  //     .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  // }
}
