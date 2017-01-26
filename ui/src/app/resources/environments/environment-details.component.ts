import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Environment } from './environment.component';
import { EnvironmentService } from './environment.service';

@Component({
  selector:    'environment',
  templateUrl: './environment-details.html',
  providers:   [EnvironmentService],
})

export class EnvironmentDetailsComponent implements OnInit {
  private editForm: FormGroup;
  private editing = false;
  private environment: Environment;

  constructor(private environmentService: EnvironmentService, private route: ActivatedRoute, private formBuilder: FormBuilder) {}

  ngOnInit() {
    this.editForm = this.formBuilder.group({
      name: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
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
}
