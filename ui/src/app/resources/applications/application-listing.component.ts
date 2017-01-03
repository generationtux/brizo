import { Observable } from 'rxjs/Rx';
import { Component, OnInit } from '@angular/core';

import { Application } from './application-details.component'
import { ApplicationService } from './application.service'

@Component({
  selector: 'application-listing',
  templateUrl: './application-listing.component.html',
  providers: [ApplicationService],
})

export class ApplicationListingComponent implements OnInit {
  private applications: Observable<Application[]>

  constructor(private applicationService: ApplicationService) {}

  ngOnInit() {
    this.applications = this.applicationService.getApplications()
  }

  logError(err: any) {
    console.error('There was an error: ' + err);
  }
}
