import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { ActivatedRoute, Params } from '@angular/router'
import { Component, EventEmitter, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Environment } from '../../environments/environment.model';
import { EnvironmentService } from '../../environments/environment.service';
import { Version } from '../version.model';
import { VersionService } from '../version.service';

@Component({
  selector:    'version',
  templateUrl: './version-details.html',
  providers:   [ EnvironmentService, VersionService ],
})

export class VersionDetailsComponent implements OnInit {
  private editForm: FormGroup;
  private editing = false;
  private version: Version;

  constructor(private versionService: VersionService,
              private route: ActivatedRoute,
              private formBuilder: FormBuilder) {}

  ngOnInit() {
    this.editForm = this.formBuilder.group({
      name: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
    });

    this.route.params
      .switchMap(
        (params: Params) => this.versionService.getVersion(params['environment-uuid'], params['version-uuid']),
      ).subscribe(
        data => {
          this.version = data;
          this.resetVersionForm(this.version);
        },
        err => console.error('There was an error: ' + err)
      );
  }

  private resetVersionForm(version: Version) {
    (<FormGroup>this.editForm)
      .setValue({
        name: version.name
      }, { onlySelf: true });
  }

  private toggleEditing() {
    this.editing = !this.editing;
  }

  private cancelEditing() {
    this.resetVersionForm(this.version);
    this.editing = false;
  }

  private saveVersion() {
    const form = (<FormGroup>this.editForm).value;
    const version = new Version(
      this.version.id,
      this.version.uuid,
      form.name,
      form.name,
      this.version.image,
      this.version.replicas,
      this.version.environment_id,
    );
    this.versionService.updateVersion(this.route.params['environment-uuid'], version.uuid, version.name)
      .subscribe(
        data => {
          this.version = data;
          this.cancelEditing()
        },
        err => console.error('There was an error: ' + err)
      );
  }
}
