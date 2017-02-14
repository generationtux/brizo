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
  }

  private handleEnvironmentData(data: any) {
    this.environment = data;
    this.application = data.application || {};
  }
}
