import { Component, OnInit } from '@angular/core';
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
export class DashboardComponent implements OnInit {
    private applications: Application[]
    private application: Application
    private applicationsIndexUrl = '/api/v1/applications'

    constructor(private applicationService: ApplicationService) {}

    ngOnInit() {
      this.applicationService.getApplications().subscribe(
        data => this.applications = data,
        err => this.logError(err)
      )
    }

    logError(err: any) {
      console.error('There was an error: ' + err);
    }
}
