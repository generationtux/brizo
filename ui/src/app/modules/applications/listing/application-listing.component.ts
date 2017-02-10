import { Observable } from 'rxjs/Rx';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { Application } from '../application.model';
import { ApplicationService } from '../application.service';

@Component({
  selector:    'application-listing',
  templateUrl: './application-listing.html',
})

export class ApplicationListingComponent implements OnInit {
  private applications: Observable<Application[]>
  private createAppForm: FormGroup;
  private modalActions: EventEmitter<string>;

  constructor(private applicationService: ApplicationService, private _fb: FormBuilder, private router: Router) {}

  ngOnInit() {
    this.applications = this.applicationService.getApplications()
    this.createAppForm = new FormGroup({
        name: new FormControl('', [<any>Validators.required])
    });
    this.modalActions = new EventEmitter<string>();
  }

  createNewApp() {
    return this.applicationService.createApplication(this.createAppForm.controls['name'].value).subscribe(
      response => this.onComplete(response),
      err => console.error(err),
    )
  }

  onComplete(response: any) {
    (<any>$('#create-application-modal')).modal('hide');
    this.router.navigate(['applications', response.uuid]);
  }
}
