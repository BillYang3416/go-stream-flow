import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FileFlowComponent } from './file-flow.component';

describe('FileFlowComponent', () => {
  let component: FileFlowComponent;
  let fixture: ComponentFixture<FileFlowComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FileFlowComponent]
    });
    fixture = TestBed.createComponent(FileFlowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
