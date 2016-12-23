import { Observable } from 'rxjs/Rx';
import { Component } from '@angular/core';

import { Application } from './application-details.component'
import { ApplicationService } from './application.service'

@Component({
  selector: 'application-creation-form',
  templateUrl: './application-creation.component.html',
  providers: [ApplicationService],
})

export class ApplicationCreationComponent {
  private things: any
  constructor(private applicationService: ApplicationService) {}

  createApplication(name: string) {
    this.applicationService.createApplication(name).subscribe(
      data => this.things = data,
      err => console.error('There was an error: ' + err),
      () => console.log('application created'),
    )

    console.log(this.things);
    return;
  }

  logError(err: any) {
    console.error('There was an error: ' + err);
  }
}
