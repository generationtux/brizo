import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { Router, ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

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
    });

    this.versionService.createVersion(version).subscribe(
      () => (this.onCreateVersion()),
      err => console.error('There was an error: ' + err)
    );
  }

  private onCreateVersion() {
    this.router.navigate(['/environments', this.environment.uuid]);
  }
}
