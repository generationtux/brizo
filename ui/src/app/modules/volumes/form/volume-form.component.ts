import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, FormControl, FormBuilder, Validators } from '@angular/forms';

import { Volume } from '../volume.model';

@Component({
    selector: 'volume-form',
    templateUrl: 'volume-form.html',
})
export class VolumeFormComponent implements OnInit {
    @Input('group')
    public volumeForm: FormGroup;

    constructor(private formBuilder: FormBuilder) {}

    ngOnInit() {
      this.volumeForm = this.formBuilder.group({
        name:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
        type:   ['', [<any>Validators.required, <any>Validators.minLength(1)]],
        source: ['', [<any>Validators.required, <any>Validators.minLength(1)]],
      })
    }
}
