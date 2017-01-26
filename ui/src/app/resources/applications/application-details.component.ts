import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params }                          from '@angular/router'
import { Component, EventEmitter, OnInit }                 from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/switchMap';

import { ApplicationService } from './application.service'
import { EnvironmentService } from '../environments/environment.service'
import { Pod }                from '../pod.component'

@Component({
  selector:    'application',
  templateUrl: './application-details.html',
  providers:   [ApplicationService],
})

export class ApplicationComponent implements OnInit {
  private application: Application;
  public editApplicationForm: FormGroup;
  public createEnvironmentForm: FormGroup;
  private environments: Environment[];
  public submitted: boolean;

  constructor(private applicationService: ApplicationService, private route: ActivatedRoute, private http: Http, private environmentService: EnvironmentService) {}

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
      uuid: new FormControl('', Validators.required)
    });
    this.createEnvironmentForm = new FormGroup({
      name: new FormControl('', Validators.required)
    });
  }

  save() {
    return this.environmentService.createEnvironment(this.createEnvironmentForm.controls['name'].value, this.application.id).subscribe(
      err => console.error('There was an error: ' + err),
      () => (this._complete()),
    )
  }
  
  update() {
    console.log(this.editApplicationForm.controls);
    return this.applicationService.editApplication(this.editApplicationForm.controls['uuid'].value, this.editApplicationForm.controls['name'].value).subscribe(
      err => console.error('There was an error: ' + err),
      () => (this._complete()),
    )
  }

  _complete() {
    (<any>$('#environment-create-modal')).modal('hide');
    (<any>$('#appliction-edit-modal')).modal('hide');
  }
}

export class Application {
  constructor(
    public id:   number,
    public ID:   number,
    public uuid: string,
    public name: string,
    public environments: Environment[],
    public pods: Array<Pod>
  ){}
}

export class Environment {
  constructor(
    public id:             number,
    public uuid:           string,
    public application_id: number,
  ){}
}
