import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';
import { Http, Response, Headers, RequestOptions } from '@angular/http';

import { Application } from '../application.model';
import { ApplicationService } from '../application.service'
import { Environment } from '../../environments/environment.model';
import { EnvironmentService } from '../../environments/environment.service';
import { Pod } from '../../pod.model';

@Component({
  selector:    'application-details',
  templateUrl: './application-details.html',
})

export class ApplicationDetailsComponent implements OnInit {
  private application: Application;
  public editApplicationForm: FormGroup;
  public createEnvironmentForm: FormGroup;
  private environments: Environment[];
  public submitted: boolean;

  constructor(
    private applicationService: ApplicationService,
    private route: ActivatedRoute,
    private router: Router,
    private http: Http,
    private environmentService: EnvironmentService
  ) {}

  ngOnInit() {
    this.route.params
      .switchMap(
        (params: Params) => this.applicationService.getApplication(params['uuid'])
      ).subscribe(
        data => {
          this.application = data
          this.environments = this.application.environments
        },
        err => console.error('There was an error: ' + err)
      );
    this.editApplicationForm = new FormGroup({
      name: new FormControl('', Validators.required),
    });
    this.createEnvironmentForm = new FormGroup({
      name: new FormControl('', Validators.required)
    });
  }

  createEnvironment() {
    return this.environmentService.createEnvironment(this.createEnvironmentForm.controls['name'].value, this.application.id).subscribe(
      (response) => (this.onCreateEnvironment(response)),
      err => console.error('There was an error: ' + err),
    )
  }

  onCreateEnvironment(response: any) {
    (<any>$('#environment-create-modal')).modal('hide');
    this.router.navigate(['environments', response.uuid]);
  }

  updateApp() {
    const application = new Application(
      this.application.id,
      this.application.id,
      this.application.uuid,
      this.editApplicationForm.controls['name'].value,
      this.application.environments,
      this.application.pods,
    );

    this.applicationService.editApplication(application).subscribe(
      (response: any) => (this.onUpdateApp(response)),
      err => console.error('There was an error: ' + err),
    )
  }

  onUpdateApp(response: any) {
    (<any>$('#appliction-edit-modal')).modal('hide');
    this.application = response
  }

}
