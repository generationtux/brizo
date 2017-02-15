import { Component, Input } from '@angular/core'
import { FormGroup } from '@angular/forms';

@Component({
    selector: 'volume',
    templateUrl: 'volume-form.html',
})
export class VolumeFormComponent {
    @Input('group')
    public volumeForm: FormGroup;
}
