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
  selector:    'version-create',
  templateUrl: './version-create.html',
})

export class VersionCreateComponent {
}
