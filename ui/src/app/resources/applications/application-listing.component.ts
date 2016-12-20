import { Observable } from 'rxjs/Rx';
import { Component, OnInit } from '@angular/core';

import { Application } from './application-details.component'
import { ApplicationService } from './application.service'

@Component({
  selector: 'application-listing',
  template: '<span>example</span>',
  providers: [ApplicationService],
})

export class ApplicationListingComponent implements OnInit {
  private applications: Application[]

  constructor(private applicationService: ApplicationService) {}

  ngOnInit() {
    this.applicationService.getApplications().subscribe(
      data => this.applications = data,
      err => this.logError(err),
      () => console.log('applications loaded')
    )
  }

  logError(err: any) {
    console.error('There was an error: ' + err);
  }
}
