import 'rxjs/add/operator/switchMap';
import { Observable } from 'rxjs/Rx';
import { Component, OnInit } from '@angular/core'
import { ActivatedRoute, Params } from '@angular/router'
import { FormArray, FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Volume } from '../volume.model';
import { VolumeService } from '../volume.service';

@Component({
    selector: 'volume-form-group',
    templateUrl: 'volume-form-group.html',
    providers:   [ VolumeService ],
})
export class VolumeFormGroupComponent implements OnInit {
    public volumesForm: FormGroup;
    public volumes: Volume[];

    constructor(private formBuilder: FormBuilder,
                private route: ActivatedRoute,
                private volumeService: VolumeService) {}

    ngOnInit() {
      this.volumesForm = this.formBuilder.group({
        volumes: this.formBuilder.array([]),
      });

      this.route.params
        .switchMap(
          (params: Params) => this.volumeService.getVolumes(params['environment-uuid'], params['version-uuid'])
        ).subscribe(
          data => this.volumes = data,
          err => console.error('There was an error: ' + err),
        )
    }

    initVolume() {
      return this.formBuilder.group({
        name:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
        type:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
        source: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
      });
    }

    addVolume() {
      const control = <FormArray>this.volumesForm.controls['volumes'];
      control.push(this.initVolume());
    }

    removeVolume(i: number) {
      const control = <FormArray>this.volumesForm.controls['volumes'];
      control.removeAt(i);
    }

    save(data: any) {
      console.log(this.volumesForm.value);
      console.log(data);
    }
}
