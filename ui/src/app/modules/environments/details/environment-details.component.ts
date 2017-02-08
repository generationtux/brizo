import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Environment } from '../environment.model';
import { EnvironmentService } from '../environment.service';
import { Version } from '../../versions/version.model';
import { VersionService } from '../../versions/version.service';

@Component({
  selector:    'environment',
  templateUrl: './environment-details.html',
  providers:   [ EnvironmentService, VersionService ],
})

export class EnvironmentDetailsComponent implements OnInit {
  private editForm: FormGroup;
  private createVersionForm: FormGroup;
  private editing = false;
  private environment: Environment;
  private versions: Observable<Version[]>

  constructor(private environmentService: EnvironmentService,
              private versionService: VersionService,
              private route: ActivatedRoute,
              private formBuilder: FormBuilder) {}

  ngOnInit() {
    this.editForm = this.formBuilder.group({
      name: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
    });

    this.createVersionForm = this.formBuilder.group({
      name: new FormControl('', Validators.required),
      image: new FormControl('', Validators.required),
      replicas: new FormControl(0, Validators.required),
    });

    this.route.params
      .switchMap(
        (params: Params) => this.environmentService.getEnvironment(params['uuid'])
      ).subscribe(
        data => {
          this.environment = data;
          this.resetEnvironmentForm(this.environment);
        },
        err => console.error('There was an error: ' + err)
      );
      
      this.versions = this.versionService.getVersions('0b0be6f1-feaf-4553-806b-eacc234dff3c');
  }

  private resetEnvironmentForm(environment: Environment) {
    (<FormGroup>this.editForm)
      .setValue({
        name: environment.name
      }, { onlySelf: true });
  }

  private toggleEditing() {
    this.editing = !this.editing;
  }

  private cancelEditing() {
    this.resetEnvironmentForm(this.environment);
    this.editing = false;
  }

  private saveEnvironment() {
    const form = (<FormGroup>this.editForm).value;
    const environment = new Environment(
      this.environment.id,
      form.name,
      form.name,
      this.environment.uuid,
      this.environment.application_id,
    );
    this.environmentService.updateEnvironment(environment)
      .subscribe(
        data => {
          this.environment = data;
          this.cancelEditing()
        },
        err => console.error('There was an error: ' + err)
      );
  }
  
  private createVersion() {
    return this.versionService.createVersion(this.createVersionForm.controls['name'].value,
                                             this.createVersionForm.controls['image'].value,
                                             parseInt(this.createVersionForm.controls['replicas'].value),
                                             this.environment.id,
                                             this.environment.uuid).subscribe(
     err => console.error('There was an error: ' + err),
     () => ((<any>$('#environment-create-modal')).modal('hide')),
    )
  }
}
