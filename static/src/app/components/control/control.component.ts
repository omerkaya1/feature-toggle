import {AfterViewInit, Component, OnChanges, OnInit, ViewChild} from '@angular/core';
import {RestService} from '../../services/rest.service';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {MatTableDataSource} from '@angular/material/table';
import {MatPaginator} from '@angular/material/paginator';

export interface Feature {
  displayName?: string;
  technicalName: string;
  expiresOn?: Date;
  description?: string;
  editMode?: boolean;
  inverted: boolean;
  active: boolean;
  customerIDs: string[];
}

@Component({
  selector: 'app-ft-control',
  templateUrl: './control.component.html',
  styleUrls: ['./control.component.css']
})
export class ControlComponent implements OnInit, OnChanges {
  public featureFG: FormGroup;
  displayedColumns: string[] = ['displayName', 'technicalName', 'expiresOn', 'description', 'inverted', 'active',
    'customerIDs', 'actions'];
  features: Feature[];
  dataSource = new MatTableDataSource<Feature>(this.features);
  editMode: boolean;

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(
    private rest: RestService,
    private fb: FormBuilder
  ) {
  }

  ngOnInit() {
    this.featureFG = this.fb.group({
      technicalName: ['', Validators.required],
      displayName: [''],
      description: [''],
      expiresOn: [],
      customerIDs: [[] as string[], Validators.required],
      inverted: [false],
      active: [false],
    });
    this.dataSource.paginator = this.paginator;
    this.editMode = false;
    this.refresh();
  }

  private refresh() {
    this.rest.getFeatures().subscribe((result: Feature[]) => {
      this.dataSource.data = result;
    });
  }

  ngOnChanges() {
    this.dataSource.paginator = this.paginator;
    this.refresh();
  }

  public editFeature(f: Feature) {
    this.rest.updateFeature(f.technicalName, f.displayName, f.description, f.expiresOn, f.active).subscribe( () => {
      f.editMode = false;
      this.refresh();
    });
  }

  public deleteFeature(name: string) {
    this.rest.deleteFeature(name).subscribe((result: string) => {
      if (result === 'feature was deleted') {
        this.refresh();
      } else {
        // TODO: replace with a proper error service
        console.log(result);
      }
    });
  }

  public submit() {
    if (typeof(this.featureFG.get('customerIDs').value) === 'string') {
      this.featureFG.get('customerIDs').setValue(this.featureFG.get('customerIDs').value.split(','));
    }
    this.rest.createFeature(this.featureFG.value).subscribe((result) => {
      if (result) {
        this.refresh();
        this.featureFG.reset();
      } else {
        // TODO: replace with a proper error service
        console.log(result);
      }
    });
  }
}
