import { Component } from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Rx';

import { Application } from '../resources/application.component'
import { ApplicationService } from '../resources/application.service'

@Component({
    selector:    'dashboard',
    templateUrl: './dashboard.component.html',
    providers: [ApplicationService],
})

@Injectable()
export class DashboardComponent {
    constructor(private applicationService: ApplicationService) {
      applicationService.getApplications().subscribe(
        data => this.applications = data,
        err => this.logError(err),
        () => console.log('applications loaded')
      )

      applicationService.getApplication('foobar').subscribe(
        data => this.application = data,
        err => this.logError(err),
        () => console.log('application loaded')
      )
    }

    private applications: Application[]

    private application: Application

    private applicationsIndexUrl = '/api/v1/applications'

    logError(err: any) {
      console.error('There was an error: ' + err);
    }
}
