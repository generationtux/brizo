import { Observable } from 'rxjs/Rx';
import { Component } from '@angular/core';

import { Application } from '../application.model';
import { ApplicationService } from '../application.service';

@Component({
  selector:    'application-creation-form',
  templateUrl: './application-creation.html',
  providers:   [ ApplicationService ],
})

export class ApplicationCreationComponent {
  constructor(private applicationService: ApplicationService) {}

  createApplication(name: string) {
    this.applicationService.createApplication(name).subscribe(
      err => console.error('There was an error: ' + err),
      () => console.log('application created'),
      () => (this._complete())
    )
  }

  logError(err: any) {
    console.error('There was an error: ' + err);
  }

  _complete() {
    (<any>$('#create-application-modal')).modal('hide');
  }
}
