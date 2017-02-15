import 'rxjs/add/operator/switchMap'
import { Observable } from 'rxjs/Rx'
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit, AfterViewInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms'

import { EnvironmentService } from '../environment.service'
import { ConfigurationService } from './configuration.service'

@Component({
  selector: 'configuration',
  templateUrl: './configuration.html',
  providers: [ConfigurationService]
})

export class ConfigurationComponent {
  private environment: any = {};
  private application: any = {};
  private createConfigForm: FormGroup;
  
  constructor(private configService: ConfigurationService,
              private environmentService: EnvironmentService,
              private route: ActivatedRoute) {}
  
  ngOnInit() {
    this.createConfigForm = new FormGroup({
        name: new FormControl('', [<any>Validators.required]),
        value: new FormControl('', [<any>Validators.required])
    });
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
  
  ngAfterViewInit() {
    var i = 1;
    $(".add-btn").click(function() {
      $('.addr' + i).html("<td><input name='key" + i + "' type='text' placeholder='KEY' class='form-control' /></td><td><input name='value" + i + "' type='text' placeholder='VALUE' class='form-control' /></td>");
      $('.config-table').append("<tr class='addr" + (i + 1) + "'></tr>");
      i++;
    });
    $(".del-btn").click(function() {
      if(i > 1) {
        $(".addr" + (i-1)).html('');
        i--;
		  }
    });
  }
  
  createNewConfig() {
    return this.configService.createConfiguration(this.createConfigForm.controls['name'].value, this.createConfigForm.controls['value'].value, this.environment.id, this.environment.uuid).subscribe(
      response => console.log('configuration created'),
      err => console.error(err),
    )
  }
}