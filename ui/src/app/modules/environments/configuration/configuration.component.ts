import 'rxjs/add/operator/switchMap'
import { Observable } from 'rxjs/Rx'
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'

import { EnvironmentService } from '../environment.service'

@Component({
  selector: 'configuration',
  templateUrl: './configuration.html',
  providers: [ EnvironmentService ]
})

export class ConfigurationComponent {
  private environment: any = {};
  private application: any = {};
  
  constructor(private environmentService: EnvironmentService, private route: ActivatedRoute) {}
  
  ngOnInit() {
    this.route.params
      .switchMap(
        (params: Params) => this.environmentService.getEnvironment(params['uuid']),
      ).subscribe(
        data => {
          this.environment = data;
          this.application = this.environment.application;
        },
        err => console.error('There was an error: ' + err)
      );
  }
}