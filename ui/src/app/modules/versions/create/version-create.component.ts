import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { Router, ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormArray, FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Environment } from '../../environments/environment.model';
import { EnvironmentService } from '../../environments/environment.service';
import { Version } from '../version.model';
import { VersionService } from '../version.service';

@Component({
  selector:    'version-create',
  templateUrl: './version-create.html',
})

export class VersionCreateComponent implements OnInit {
  private application: any = {}
  private environment: any = {}
  private versionForm: FormGroup;

  constructor(
    private versionService: VersionService,
    private environmentService: EnvironmentService,
    private activeRoute: ActivatedRoute,
    private router: Router,
    private formBuilder: FormBuilder
  ) {}

  ngOnInit() {
    this.activeRoute.params.switchMap((params: Params) => {
      return this.environmentService.getEnvironment(params['environment-uuid']);
    }).subscribe(
      data => this.handleEnvironmentData(data),
      err => console.error('There was an error: ' + err)
    );

    this.versionForm = this.formBuilder.group({
      name: new FormControl('', Validators.required),
      image: new FormControl('', Validators.required),
      replicas: new FormControl(1, Validators.required),
      pullPolicy: new FormControl('', Validators.required),
      args: new FormControl(),
      ports: new FormControl(),
      volumes:  this.formBuilder.array([]),
    });
  }

  private handleEnvironmentData(data: any) {
    this.environment = data;
    this.application = data.application || {};
  }

  private createVersion() {
    let version = new Version({
      name: this.versionForm.controls['name'].value,
      image: this.versionForm.controls['image'].value,
      replicas: parseInt(this.versionForm.controls['replicas'].value),
      environment_uuid: this.environment.uuid,
      volumes: this.versionForm.controls['volumes'].value,
    });

    this.versionService.createVersion(version).subscribe(
      () => (this.onCreateVersion()),
      err => console.error('There was an error: ' + err)
    );
  }

  private onCreateVersion() {
    this.router.navigate(['/environments', this.environment.uuid]);
  }

  initVolume() {
    return this.formBuilder.group({
      name:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
      type:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
      source: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
    });
  }

  addVolume() {
    console.log('add')
    const control = <FormArray>this.versionForm.controls['volumes'];
    const volCtrl = this.initVolume();
    control.push(volCtrl);
  }

  removeVolume(i: number) {
    console.log('remove')
    const control = <FormArray>this.versionForm.controls['volumes'];
    control.removeAt(i);
  }
}
