import 'rxjs/add/operator/switchMap'
import { Observable } from 'rxjs/Rx'
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit, AfterViewInit } from '@angular/core'
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
  
  ngAfterViewInit() {
    var i = 1;
    console.log(i)
    $(".add-btn").click(function() {
      console.log(i)
      $('.addr' + i).html("<td><input name='key" + i + "' type='text' placeholder='KEY' class='form-control' /></td><td><input name='value" + i + "' type='text' placeholder='VALUE' class='form-control' /></td>");
      $('.config-table').append("<tr class='addr" + (i + 1) + "'></tr>");
      i++;
    });
    $(".del-btn").click(function() {
      if(i > 1) {
        $(".addr" + (i-1)).html('');
        i--;
		  }
      console.log(i)
    });
  }
  
  save() {
    
  }
}