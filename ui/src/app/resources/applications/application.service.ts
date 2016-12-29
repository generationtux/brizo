import { Observable } from 'rxjs/Rx';
import { Component, Injectable } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { Application } from './application-details.component'

@Injectable()
export class ApplicationService {

  private applicationsCreateUrl = '/api/v1/applications'
  private applicationsGetUrl    = '/api/v1/applications/'
  private applicationsIndexUrl  = '/api/v1/applications'

  constructor(private http: Http) { }

  getApplications(): Observable<Application[]> {
    return this.http.get(this.applicationsIndexUrl)
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  getApplication(uuid: string): Observable<Application> {
    return this.http.get(this.applicationsGetUrl + `${uuid}`)
      .map((res: Response) => res.json())
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }

  createApplication(name: string): Observable<Application> {
    const headers = new Headers({ 'Content-Type': 'application/json' });
    const options = new RequestOptions({ headers: headers });

    return this.http.post(this.applicationsCreateUrl, { name }, options)
      .map((res: Response) => res.json() || {})
      .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
  }
}
