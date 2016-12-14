import { Component } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';
import { Application } from '../resources/application.component'

@Component({
    selector:    'dashboard',
    templateUrl: './dashboard.component.html',
})

@Injectable()
export class DashboardComponent {
    constructor(private http: Http) {
      this.getApplications().subscribe(
        data => this.applications = data,
        err => this.logError(err),
        () => console.log('applications loaded')
      )
    }

    private applications: Application[]

    private applicationsIndexUrl = '/api/v1/applications'

    logError(err: any) {
      console.error('There was an error: ' + err);
    }

    getApplications(): Observable<Application[]> {
      return this.http.get(this.applicationsIndexUrl)
        .map((res: Response) => res.json())
        .catch((error: any) => Observable.throw(error.json().error || 'Server error'));
    }
}
